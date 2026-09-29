import { useEffect, useMemo, useRef, useState } from 'react';
import { Image, Linking, Modal, PanResponder, Platform, Pressable, ScrollView, StyleSheet, Text, TextInput, View } from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';

export type SearchArea = { latitude: number; longitude: number; label: string; device?: boolean };
type Place = { lat: string; lon: string; display_name: string };
type Props = {
  area: SearchArea;
  locating: boolean;
  locationError: string | null;
  onLocate: () => Promise<{ latitude: number; longitude: number } | undefined>;
  onClose: () => void;
  onApply: (area: SearchArea) => void;
};

const TILE = 256;
function project(area: SearchArea, zoom: number) {
  const size = TILE * 2 ** zoom;
  const latitude = Math.max(-85.0511, Math.min(85.0511, area.latitude)) * Math.PI / 180;
  return { x: (area.longitude + 180) / 360 * size, y: (1 - Math.asinh(Math.tan(latitude)) / Math.PI) / 2 * size };
}
function unproject(x: number, y: number, zoom: number): SearchArea {
  const size = TILE * 2 ** zoom;
  return {
    longitude: ((x / size * 360) % 360 + 360) % 360 - 180,
    latitude: Math.atan(Math.sinh(Math.PI * (1 - 2 * Math.max(0, Math.min(size, y)) / size))) * 180 / Math.PI,
    label: 'Selected area',
  };
}

function Action({ label, accessibilityLabel, onPress, disabled = false }: { label: string; accessibilityLabel?: string; onPress: () => void; disabled?: boolean }) {
  return <Pressable accessibilityRole="button" accessibilityLabel={accessibilityLabel ?? label} disabled={disabled} accessibilityState={{ disabled }} onPress={onPress} style={({ pressed }) => [styles.action, (pressed || disabled) && { opacity: 0.5 }]}><Text style={styles.actionText}>{label}</Text></Pressable>;
}

export function LocationSheet({ area, locating, locationError, onLocate, onClose, onApply }: Props) {
  const [draft, setDraft] = useState(area);
  const [zoom, setZoom] = useState(15);
  const [size, setSize] = useState({ width: 0, height: 280 });
  const [query, setQuery] = useState('');
  const [places, setPlaces] = useState<Place[]>([]);
  const [searching, setSearching] = useState(false);
  const [searchError, setSearchError] = useState<string | null>(null);
  const [mapError, setMapError] = useState(false);
  const alive = useRef(true);
  const searchRequest = useRef<AbortController | null>(null);
  const lastSearch = useRef(0);
  const insets = useSafeAreaInsets();
  const center = project(draft, zoom);
  const current = useRef({ center, zoom });
  current.current = { center, zoom };
  const gesture = useRef({ x: 0, y: 0, zoom: 15 });
  useEffect(() => {
    alive.current = true;
    return () => { alive.current = false; searchRequest.current?.abort(); };
  }, []);

  const pan = useMemo(() => PanResponder.create({
    onMoveShouldSetPanResponder: (_, g) => Math.abs(g.dx) + Math.abs(g.dy) > 4,
    onPanResponderGrant: () => { gesture.current = { ...current.current.center, zoom: current.current.zoom }; },
    onPanResponderMove: (_, g) => setDraft(unproject(gesture.current.x - g.dx, gesture.current.y - g.dy, gesture.current.zoom)),
  }), []);

  async function search() {
    if (!query.trim() || searching || Date.now() - lastSearch.current < 1100) return;
    lastSearch.current = Date.now();
    searchRequest.current?.abort();
    const controller = new AbortController();
    searchRequest.current = controller;
    const timeout = setTimeout(() => controller.abort(), 10000);
    setSearching(true);
    setPlaces([]);
    setSearchError(null);
    try {
      const viewbox = [Math.max(-180, draft.longitude - 0.5), Math.min(85, draft.latitude + 0.5), Math.min(180, draft.longitude + 0.5), Math.max(-85, draft.latitude - 0.5)].join(',');
      const response = await fetch(`https://nominatim.openstreetmap.org/search?format=jsonv2&limit=4&viewbox=${viewbox}&q=${encodeURIComponent(query.trim())}`, {
        signal: controller.signal,
        headers: Platform.OS === 'web' ? undefined : { 'User-Agent': 'CheapEats/1.0' },
      });
      if (!response.ok) throw new Error('Search failed');
      const data: Place[] = await response.json();
      if (!alive.current) return;
      const valid = data.filter(place => Number.isFinite(Number(place.lat)) && Number.isFinite(Number(place.lon)) && Math.abs(Number(place.lat)) <= 85.0511 && Math.abs(Number(place.lon)) <= 180);
      setPlaces(valid);
      if (!valid.length) setSearchError('No places found. Try a station, address, or neighborhood.');
    } catch {
      if (alive.current) setSearchError('Location search is unavailable. Try again or move the map.');
    } finally {
      clearTimeout(timeout);
      if (alive.current) setSearching(false);
    }
  }

  const tiles = [];
  const left = center.x - size.width / 2;
  const top = center.y - size.height / 2;
  const count = 2 ** zoom;
  if (size.width) for (let x = Math.floor(left / TILE); x <= Math.floor((left + size.width) / TILE); x++) {
    for (let y = Math.floor(top / TILE); y <= Math.floor((top + size.height) / TILE); y++) {
      if (y < 0 || y >= count) continue;
      tiles.push(<Image key={`${zoom}/${x}/${y}`} source={{ uri: `https://tile.openstreetmap.org/${zoom}/${((x % count) + count) % count}/${y}.png`, ...(Platform.OS !== 'web' ? { headers: { 'User-Agent': 'CheapEats/1.0' } } : {}) }} onError={() => setMapError(true)} style={{ position: 'absolute', left: x * TILE - left, top: y * TILE - top, width: TILE, height: TILE }} />);
    }
  }
  function nudge(dx: number, dy: number) { setDraft(unproject(center.x + dx, center.y + dy, zoom)); }

  return <Modal visible transparent animationType="slide" onRequestClose={onClose}>
    <View style={styles.overlay}>
      <Pressable accessibilityRole="button" accessibilityLabel="Close location picker" onPress={onClose} style={styles.backdrop} />
      <View accessibilityViewIsModal onAccessibilityEscape={onClose} style={[styles.sheet, { paddingBottom: Math.max(insets.bottom, 20) }]}>
        <ScrollView keyboardShouldPersistTaps="handled" contentContainerStyle={styles.content}>
          <View style={styles.row}><Text accessibilityRole="header" style={styles.title}>Choose location</Text><Action label="Close" onPress={onClose} /></View>
          <View style={styles.row}>
            <TextInput accessibilityLabel="Search station, address, or neighborhood" placeholder="Station, address, or neighborhood" value={query} onChangeText={setQuery} onSubmitEditing={() => void search()} returnKeyType="search" style={styles.input} />
            <Action label={searching ? 'Searching…' : 'Search'} disabled={searching || !query.trim()} onPress={() => void search()} />
          </View>
          {searchError && <Text accessibilityRole="alert" style={styles.hint}>{searchError}</Text>}
          {places.map(place => <Pressable key={`${place.lat}/${place.lon}`} accessibilityRole="button" style={styles.result} onPress={() => { setDraft({ latitude: Number(place.lat), longitude: Number(place.lon), label: place.display_name.split(',')[0] }); setPlaces([]); setZoom(15); setMapError(false); }}><Text style={styles.actionText}>{place.display_name}</Text></Pressable>)}
          <View>
            <Action label={locating ? 'Finding location…' : 'Use my current location'} disabled={locating} onPress={() => { void onLocate().then(position => { if (position && alive.current) { setDraft({ ...position, label: 'My location', device: true }); setZoom(15); setPlaces([]); setMapError(false); } }); }} />
            <Text style={styles.hint}>Your coordinates are sent to Cheap Eats when you search this area. Map and place searches use OpenStreetMap.</Text>
            {locationError && <Text accessibilityRole="alert" style={styles.hint}>{locationError}</Text>}
          </View>
          <View style={[styles.map, Platform.OS === 'web' && { touchAction: 'none' }]} onLayout={event => setSize(event.nativeEvent.layout)} {...pan.panHandlers}>
            <View style={[StyleSheet.absoluteFill, { pointerEvents: 'none' }]}>
              {tiles}
            </View>

            <View style={[styles.pin, { pointerEvents: 'none' }]}>
              <Text style={styles.pinText}>📍</Text>
            </View>
            <View style={styles.zoom}><Action label="+" accessibilityLabel="Zoom in" onPress={() => setZoom(Math.min(18, zoom + 1))} disabled={zoom === 18} /><Action label="−" accessibilityLabel="Zoom out" onPress={() => setZoom(Math.max(3, zoom - 1))} disabled={zoom === 3} /></View>
          </View>
          <View style={styles.row}><Text style={styles.hint}>Move the map to place the pin</Text><Pressable accessibilityRole="link" onPress={() => void Linking.openURL('https://www.openstreetmap.org/copyright')}><Text style={styles.attribution}>© OpenStreetMap contributors</Text></Pressable></View>
          {mapError && <Text accessibilityRole="alert" style={styles.hint}>Some map tiles could not load. Check your connection or search for a place above.</Text>}
          <View style={styles.row}><Text numberOfLines={2} style={[styles.actionText, { flex: 1 }]}>{draft.label}</Text><View style={styles.row}><Action label="←" accessibilityLabel="Move map west" onPress={() => nudge(-100, 0)} /><Action label="↑" accessibilityLabel="Move map north" onPress={() => nudge(0, -100)} /><Action label="↓" accessibilityLabel="Move map south" onPress={() => nudge(0, 100)} /><Action label="→" accessibilityLabel="Move map east" onPress={() => nudge(100, 0)} /></View></View>
          <Pressable accessibilityRole="button" disabled={locating} accessibilityState={{ disabled: locating }} onPress={() => onApply(draft)} style={[styles.apply, locating && { opacity: 0.5 }]}><Text style={styles.applyText}>Search this area</Text></Pressable>
        </ScrollView>
      </View>
    </View>
  </Modal>;
}

const styles = StyleSheet.create({
  overlay: { flex: 1, justifyContent: 'flex-end' },
  backdrop: { ...StyleSheet.absoluteFill, backgroundColor: 'rgba(0,0,0,0.35)' },
  sheet: { width: '100%', maxWidth: 580, maxHeight: '92%', alignSelf: 'center', backgroundColor: '#FFFFFF', borderTopLeftRadius: 24, borderTopRightRadius: 24, overflow: 'hidden' },
  content: { padding: 20, gap: 12 },
  row: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', gap: 4, flexWrap: 'wrap' },
  title: { fontSize: 23, fontWeight: '700', color: '#171717' },
  action: { minWidth: 44, minHeight: 44, paddingHorizontal: 8, alignItems: 'center', justifyContent: 'center' },
  actionText: { fontSize: 13, fontWeight: '600', color: '#171717' },
  input: { flex: 1, minWidth: 160, minHeight: 44, borderWidth: 0.5, borderColor: '#888888', borderRadius: 12, paddingHorizontal: 12, fontSize: 13, color: '#171717' },
  hint: { fontSize: 12, lineHeight: 18, color: '#707070' },
  result: { paddingVertical: 12, borderBottomWidth: 0.5, borderBottomColor: '#DDDDDD' },
  map: { height: 280, overflow: 'hidden', backgroundColor: '#EEEEEE', borderRadius: 12 },
  pin: { position: 'absolute', top: '50%', left: '50%', marginLeft: -18, marginTop: -36 },
  pinText: { fontSize: 36, lineHeight: 40 },
  zoom: { position: 'absolute', right: 8, top: 8, backgroundColor: '#FFFFFF', borderRadius: 10 },
  attribution: { fontSize: 10, color: '#555555', textDecorationLine: 'underline', paddingVertical: 12 },
  apply: { minHeight: 48, justifyContent: 'center', alignItems: 'center', borderRadius: 24, backgroundColor: '#171717' },
  applyText: { fontSize: 14, fontWeight: '600', color: '#FFFFFF' },
});

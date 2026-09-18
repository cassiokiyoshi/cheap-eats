import { useRef, useState } from 'react';
import { FlatList, Pressable, StyleSheet, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

type Dish = { id: number; japanese: string; english: string; restaurant: string; price: number; distance: number };
// Fictional preview data. No API or device location is used yet.
const dishes: Dish[] = [
  { id: 1, japanese: 'たこ焼き 8個', english: 'Takoyaki · 8 pcs', restaurant: 'Takoyaki Taro', price: 500, distance: 90 },
  { id: 2, japanese: '醤油ラーメン', english: 'Shoyu Ramen', restaurant: 'Yamada Noodles', price: 850, distance: 150 },
  { id: 3, japanese: '味噌ラーメン', english: 'Miso Ramen', restaurant: 'Yamada Noodles', price: 900, distance: 150 },
  { id: 4, japanese: 'チャーシュー丼', english: 'Chashu Rice Bowl', restaurant: 'Bowl Kitchen', price: 600, distance: 180 },
  { id: 5, japanese: '餃子 6個', english: 'Gyoza · 6 pcs', restaurant: 'Gyoza House', price: 450, distance: 220 },
  { id: 6, japanese: 'ビーフカレー', english: 'Beef Curry Rice', restaurant: 'Curry Corner', price: 800, distance: 250 },
  { id: 7, japanese: 'かけうどん', english: 'Udon Noodles', restaurant: 'Udon Corner', price: 550, distance: 260 },
  { id: 8, japanese: '親子丼', english: 'Chicken & Egg Bowl', restaurant: 'Bowl Kitchen', price: 780, distance: 270 },
  { id: 9, japanese: 'おにぎりセット', english: 'Onigiri Set', restaurant: 'Rice & Co.', price: 480, distance: 280 },
  { id: 10, japanese: '焼き魚定食', english: 'Grilled Fish Set', restaurant: 'Lunch House', price: 980, distance: 290 },
  { id: 11, japanese: 'ざるそば', english: 'Cold Soba', restaurant: 'Soba House', price: 700, distance: 295 },
  { id: 12, japanese: 'チキンカレー', english: 'Chicken Curry', restaurant: 'Curry Corner', price: 900, distance: 300 },
];
const yen = (value: number) => `¥${value.toLocaleString('en-US')}`;
const PAGE_SIZE = 10;

function Choice({ label, selected = false, disabled = false, onPress }: { label: string; selected?: boolean; disabled?: boolean; onPress: () => void }) {
  return <Pressable accessibilityRole="button" accessibilityState={{ selected, disabled }} disabled={disabled} onPress={onPress}
    style={({ pressed }) => [styles.button, selected && styles.selected, disabled && styles.disabled, pressed && styles.pressed]}>
    <Text style={[styles.buttonText, selected && styles.selectedText]}>{label}</Text>
  </Pressable>;
}

export default function HomeScreen() {
  const [sort, setSort] = useState<'distance' | 'price'>('distance');
  const [budget, setBudget] = useState(1000);
  const [radius, setRadius] = useState(300);
  const [filtersOpen, setFiltersOpen] = useState(false);
  const [page, setPage] = useState(0);
  const list = useRef<FlatList<Dish>>(null);
  const results = dishes.filter((dish) => dish.price <= budget && dish.distance <= radius)
    .sort((a, b) => sort === 'distance'
      ? a.distance - b.distance || a.price - b.price || a.id - b.id
      : a.price - b.price || a.distance - b.distance || a.id - b.id);
  const pageCount = Math.max(1, Math.ceil(results.length / PAGE_SIZE));
  function changePage(next: number) { setPage(next); list.current?.scrollToOffset({ offset: 0, animated: true }); }

  return <SafeAreaView style={styles.screen}>
    <View style={styles.container}>
      <View style={styles.header}>
        <View style={styles.topRow}>
          <Text accessibilityRole="header" style={styles.brand}>Cheap Eats</Text>
          <Pressable accessibilityRole="button" accessibilityLabel="Filters" accessibilityState={{ expanded: filtersOpen }}
            onPress={() => setFiltersOpen(!filtersOpen)} style={[styles.button, filtersOpen && styles.filterActive]}>
            <Text style={styles.buttonText}>Filters {filtersOpen ? '⌃' : '⌄'}</Text>
          </Pressable>
        </View>
        <Text style={styles.summary}>Within {radius} m · Up to {yen(budget)}</Text>
        <View style={styles.sortRow}>
          <View style={styles.options}>
            <Choice label="Distance" selected={sort === 'distance'} onPress={() => { setSort('distance'); changePage(0); }} />
            <Choice label="Price" selected={sort === 'price'} onPress={() => { setSort('price'); changePage(0); }} />
          </View>
          <Text style={styles.muted}>10 per page</Text>
        </View>
        {filtersOpen && <View style={styles.filterPanel}>
          <Text style={styles.label}>Maximum price</Text>
          <View style={styles.options}>{[500, 800, 1000, 1500].map((value) =>
            <Choice key={value} label={yen(value)} selected={budget === value} onPress={() => { setBudget(value); changePage(0); }} />)}</View>
          <Text style={styles.label}>Search radius</Text>
          <View style={styles.options}>{[100, 300, 500, 1000].map((value) =>
            <Choice key={value} label={`${value} m`} selected={radius === value} onPress={() => { setRadius(value); changePage(0); }} />)}</View>
          <View style={styles.options}>
            <Choice label="Reset" onPress={() => { setBudget(1000); setRadius(300); changePage(0); }} />
            <Choice label="Done" selected onPress={() => setFiltersOpen(false)} />
          </View>
        </View>}
      </View>
      <FlatList ref={list} data={results.slice(page * PAGE_SIZE, (page + 1) * PAGE_SIZE)} numColumns={2}
        keyExtractor={(dish) => String(dish.id)} style={styles.list} contentContainerStyle={styles.content} columnWrapperStyle={styles.columns}
        renderItem={({ item }) => <View style={styles.card}>
          <View style={styles.photo} accessibilityLabel="Dish photo placeholder">
            <Text style={styles.photoMark}>写真</Text><Text style={styles.photoLabel}>PHOTO PREVIEW</Text>
          </View>
          <View style={styles.cardBody}>
            <View style={styles.nameRow}><Text style={styles.japanese}>{item.japanese}</Text><Text style={styles.price}>{yen(item.price)}</Text></View>
            <Text style={styles.english}>{item.english}</Text>
            <View style={styles.metaRow}><Text style={styles.restaurant}>{item.restaurant}</Text><Text style={styles.distance}>{item.distance} m</Text></View>
          </View>
        </View>}
        ListEmptyComponent={<View style={styles.empty}><Text style={styles.label}>No dishes in this range</Text><Text style={styles.summary}>Try a higher budget or a wider radius.</Text></View>}
        ListFooterComponent={<View style={styles.footer}>
          <View style={styles.pagination}>
            <Choice label="Previous" disabled={page === 0} onPress={() => changePage(page - 1)} />
            <Text style={styles.muted}>Page {page + 1} of {pageCount}</Text>
            <Choice label="Next" disabled={page + 1 >= pageCount} onPress={() => changePage(page + 1)} />
          </View>
          <Text style={styles.preview}>Design preview · Sample dishes and distances</Text>
        </View>} />
    </View>
  </SafeAreaView>;
}

const styles = StyleSheet.create({
  screen: { flex: 1, backgroundColor: '#FAFAFA' },
  container: { flex: 1, width: '100%', maxWidth: 580, alignSelf: 'center' },
  header: { paddingHorizontal: 16, paddingTop: 20, paddingBottom: 16, backgroundColor: '#FFFFFF', borderBottomWidth: 0.5, borderBottomColor: '#B8B8B8' },
  topRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', gap: 12 },
  brand: { fontSize: 28, fontWeight: '800', letterSpacing: -1, color: '#151515', flexShrink: 1 },
  button: { minHeight: 44, paddingHorizontal: 14, justifyContent: 'center', borderWidth: 0.5, borderColor: '#777777', borderRadius: 22 },
  buttonText: { fontSize: 13, fontWeight: '600', color: '#222222' },
  filterActive: { backgroundColor: '#EEEEEE' },
  summary: { fontSize: 12, color: '#666666', marginTop: 8, lineHeight: 18 },
  sortRow: { marginTop: 18, flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', gap: 8 },
  selected: { backgroundColor: '#191919', borderColor: '#191919' },
  selectedText: { color: '#FFFFFF' },
  muted: { fontSize: 12, color: '#666666' },
  filterPanel: { marginTop: 16, padding: 14, backgroundColor: '#F7F7F7', borderWidth: 0.5, borderColor: '#B8B8B8', borderRadius: 14, gap: 12 },
  label: { fontSize: 14, fontWeight: '600', color: '#333333' },
  options: { flexDirection: 'row', flexWrap: 'wrap', gap: 8 },
  list: { flex: 1 },
  content: { padding: 14, paddingBottom: 32 },
  columns: { gap: 12, marginBottom: 14 },
  card: { flex: 1, maxWidth: '50%', backgroundColor: '#FFFFFF', borderRadius: 16, borderWidth: 0.5, borderColor: 'rgba(0,0,0,0.32)', overflow: 'hidden' },
  photo: { width: '100%', aspectRatio: 0.9, backgroundColor: '#EEEEEE', alignItems: 'center', justifyContent: 'center', gap: 10, borderBottomWidth: 0.5, borderBottomColor: '#D4D4D4' },
  photoMark: { fontSize: 28, color: '#A0A0A0', fontWeight: '300' },
  photoLabel: { fontSize: 8, letterSpacing: 2, color: '#777777' },
  cardBody: { padding: 10, flex: 1 },
  nameRow: { flexDirection: 'row', gap: 4, alignItems: 'flex-start' },
  japanese: { flex: 1, fontSize: 14, lineHeight: 21, fontWeight: '600', color: '#171717' },
  price: { fontSize: 16, lineHeight: 21, fontWeight: '700', color: '#171717' },
  english: { fontSize: 12, lineHeight: 18, color: '#666666', marginTop: 3 },
  metaRow: { flexDirection: 'row', flexWrap: 'wrap', justifyContent: 'space-between', alignItems: 'flex-end', gap: 5, marginTop: 20 },
  restaurant: { fontSize: 11, lineHeight: 16, color: '#666666', flexShrink: 1 },
  distance: { fontSize: 11, lineHeight: 16, color: '#666666' },
  footer: { gap: 22, marginTop: 10 },
  pagination: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', gap: 8 },
  disabled: { opacity: 0.35 },
  pressed: { opacity: 0.65 },
  preview: { textAlign: 'center', fontSize: 11, color: '#777777', lineHeight: 18 },
  empty: { paddingVertical: 50, alignItems: 'center' },
});

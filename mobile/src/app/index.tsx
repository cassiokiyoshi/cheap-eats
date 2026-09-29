import { useCallback, useEffect, useRef, useState } from 'react';
import { fetchNearbyDishes } from '@/api/dishes';
import { FlatList, Pressable, StyleSheet, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { FilterSheet } from '@/components/filter-sheet';
import { LocationSheet, type SearchArea } from '@/components/location-sheet';
import { router, useFocusEffect } from 'expo-router';
import type { ImageSourcePropType } from 'react-native';
import { DishPhoto } from '@/components/dish-photo';
import { getDemoDishImage } from '@/constants/demo-images';
import * as Location from 'expo-location';

type Dish = {
  id: number;
  japanese: string;
  english: string;
  restaurant: string;
  price: number;
  distance: number;
  image?: ImageSourcePropType;
  isDemo: boolean;
};

const yen = (value: number) => `¥${value.toLocaleString('en-US')}`;
const PAGE_SIZE = 10;
const TOKYO_STATION = {
  latitude: 35.6812,
  longitude: 139.7671,
};

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

  const [results, setResults] = useState<Dish[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hasNext, setHasNext] = useState(false);
  const [retryCount, setRetryCount] = useState(0);

  const [searchLocation, setSearchLocation] = useState(TOKYO_STATION);
  const [locationOpen, setLocationOpen] = useState(false);
  const [locationLabel, setLocationLabel] = useState('Tokyo Station');
  const [previewArea, setPreviewArea] = useState(true);
  const [usingDeviceLocation, setUsingDeviceLocation] = useState(false);
  const [locating, setLocating] = useState(false);
  const [locationError, setLocationError] = useState<string | null>(null);

  const locationRequest = useRef(0);
  const locationBusy = useRef(false);
  const locationTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  const list = useRef<FlatList<Dish>>(null);

  function changePage(next: number) {
    setPage(next);
    list.current?.scrollToOffset({ offset: 0, animated: true });
  }

  useEffect(() => {
    return () => {
      locationRequest.current += 1;
      locationBusy.current = false;

      if (locationTimer.current !== null) {
        clearTimeout(locationTimer.current);
      }
    };
  }, []);

  function closeLocation() {
    // Ignore any location request that finishes after this action.
    locationRequest.current += 1;
    locationBusy.current = false;

    if (locationTimer.current !== null) {
      clearTimeout(locationTimer.current);
      locationTimer.current = null;
    }

    setLocating(false);
    setLocationError(null);
    setLocationOpen(false);
  }

  async function useMyLocation() {
    if (locationBusy.current) return;

    locationBusy.current = true;
    const request = ++locationRequest.current;

    setLocating(true);
    setLocationError(null);

    try {
      const permission = await Location.requestForegroundPermissionsAsync();

      if (request !== locationRequest.current) return;

      if (!permission.granted) {
        setLocationError(
          permission.canAskAgain
            ? 'Location permission was declined. Your search area is unchanged.'
            : 'Enable location permission in your browser or device settings to use this feature.',
        );
        return;
      }

      // Start the timeout after the permission prompt is answered.
      locationTimer.current = setTimeout(() => {
        if (request !== locationRequest.current) return;

        locationRequest.current += 1;
        locationBusy.current = false;
        locationTimer.current = null;

        setLocating(false);
        setLocationError(
          'Finding your location took too long. Your search area is unchanged. Try again.',
        );
      }, 15000);

      const position = await Location.getCurrentPositionAsync({
        accuracy: Location.Accuracy.High,
      });

      if (request !== locationRequest.current) return;

      const { latitude, longitude } = position.coords;

      if (
        !Number.isFinite(latitude) ||
        !Number.isFinite(longitude) ||
        latitude < -90 ||
        latitude > 90 ||
        longitude < -180 ||
        longitude > 180
      ) {
        throw new Error('Invalid coordinates');
      }

      return { latitude, longitude };
    } catch {
      if (request !== locationRequest.current) return;

      setLocationError(
        'Could not read your location. Check location services and permissions. Your search area is unchanged.',
      );
    } finally {
      if (request === locationRequest.current) {
        if (locationTimer.current !== null) {
          clearTimeout(locationTimer.current);
          locationTimer.current = null;
        }

        locationBusy.current = false;
        setLocating(false);
      }
    }
  }

  useFocusEffect(useCallback(() => {
    const controller = new AbortController();
    let active = true;
    let timedOut = false;

    setLoading(true);
    setError(null);
    setResults([]);
    setHasNext(false);

    const timeout = setTimeout(() => {
      timedOut = true;
      controller.abort();
    }, 10000);

    async function loadDishes() {
      try {
        const data = await fetchNearbyDishes(
          {
            latitude: searchLocation.latitude,
            longitude: searchLocation.longitude,
            radius,
            maxPrice: budget,
            sort,
            limit: PAGE_SIZE + 1,
            offset: page * PAGE_SIZE,
          },
          controller.signal,
        );

        if (!active) return;

        const cards: Dish[] = data.slice(0, PAGE_SIZE).map((result) => {
          const dishName = result.dish.name_ja ?? result.dish.name;
          const restaurantName =
            result.restaurant.name_ja ?? result.restaurant.name;

          const restaurantEnglish = result.restaurant.name_en;

          return {
            id: result.dish.id,
            japanese: dishName,
            english:
              result.dish.name_en && result.dish.name_en !== dishName
                ? result.dish.name_en
                : '',
            restaurant: [
              restaurantName,
              restaurantEnglish !== restaurantName
                ? restaurantEnglish
                : null,
            ]
              .filter(Boolean)
              .join('\n'),
            price: result.dish.price,
            distance: result.distance_meters,
            image: getDemoDishImage(
              result.restaurant.name,
              result.dish.name_en,
            ),
            isDemo: result.restaurant.name.startsWith('[DEMO] '),
          };
        });

        setResults(cards);
        setHasNext(data.length > PAGE_SIZE);
      } catch (cause) {
        if (!active) return;

        setError(
          timedOut
            ? 'The request timed out. Please try again.'
            : cause instanceof Error
              ? cause.message
              : 'Could not load nearby dishes.',
        );
      } finally {
        clearTimeout(timeout);
        if (active) setLoading(false);
      }
    }

    void loadDishes();

    return () => {
      active = false;
      clearTimeout(timeout);
      controller.abort();
    };
  }, [
    sort,
    budget,
    radius,
    page,
    retryCount,
    searchLocation.latitude,
    searchLocation.longitude,
  ]));

  return <SafeAreaView style={styles.screen}>
    <View style={styles.container}>
      <View style={styles.header}>
        <View style={styles.topRow}>
          <Text accessibilityRole="header" style={styles.brand}>Cheap Eats</Text>
          <Pressable
            accessibilityRole="button"
            accessibilityLabel="Open filters"
            accessibilityState={{ expanded: filtersOpen }}
            onPress={() => setFiltersOpen(true)}
            style={({ pressed }) => [
              styles.filterIconButton,
              pressed && styles.pressed,
            ]}
          >
            <View accessible={false} style={styles.filterGlyph}>
              <View style={[styles.filterStripe, { width: 24 }]} />
              <View style={[styles.filterStripe, { width: 16 }]} />
              <View style={[styles.filterStripe, { width: 8 }]} />
            </View>
          </Pressable>
        </View>
          <View
            style={{
              flexDirection: 'row',
              alignItems: 'center',
              justifyContent: 'space-between',
              gap: 12,
            }}
          >
            <Pressable
              accessibilityRole="button"
              accessibilityLabel={`Change search location, ${locationLabel}${previewArea ? ', preview area' : ''}`}
              accessibilityState={{ expanded: locationOpen }}
              onPress={() => setLocationOpen(true)}
              style={({ pressed }) => [
                styles.locationSelector,
                { flex: 1, minWidth: 0 },
                pressed && styles.pressed,
              ]}
            >
              <Text style={styles.locationLink} numberOfLines={1}>
                ⌖ {locationLabel} ▾
              </Text>
              {previewArea && <Text style={styles.previewBadge}>Preview</Text>}
            </Pressable>

            <Pressable
              accessibilityRole="button"
              onPress={() => router.push('/dishes/new')}
              style={({ pressed }) => [
                {
                  minHeight: 44,
                  justifyContent: 'center',
                  flexShrink: 0,
                },
                pressed && styles.pressed,
              ]}
            >
              <Text
                style={{
                  fontSize: 12,
                  color: '#171717',
                  textDecorationLine: 'underline',
                }}
              >
                + Add dish
              </Text>
            </Pressable>
          </View>
        <View style={styles.compactControls}>
          <Text style={styles.filterSummary}>Within {radius} m · Up to {yen(budget)}</Text>
          <Pressable accessibilityRole="button" accessibilityLabel={`Sorted by ${sort === 'price' ? 'lowest price' : 'nearest'}. Switch to ${sort === 'price' ? 'nearest' : 'lowest price'}`} onPress={() => { setSort(sort === 'price' ? 'distance' : 'price'); changePage(0); }} style={styles.sortControl}>
            <Text style={styles.sortText}>Sort: {sort === 'price' ? 'Price' : 'Nearest'} ↕</Text>
          </Pressable>
        </View>
      </View>
      <FlatList ref={list} data={results} numColumns={2}
        keyExtractor={(dish) => String(dish.id)} style={styles.list} contentContainerStyle={styles.content} columnWrapperStyle={styles.columns}
        renderItem={({ item }) => (
          <Pressable
            accessibilityRole="button"
            accessibilityLabel={`View ${item.english || item.japanese}`}
            onPress={() =>
              router.push({
                pathname: '/dishes/[id]',
                params: { id: String(item.id) },
              })
            }
            style={({ pressed }) => [
              styles.card,
              pressed && styles.pressed,
            ]}
          >
            <DishPhoto
              source={item.image}
              label={item.english || item.japanese}
            />

            <View style={styles.cardBody}>
              <View style={styles.nameRow}>
                <Text style={styles.japanese}>{item.japanese}</Text>
                <Text style={styles.price}>{yen(item.price)}</Text>
              </View>

              {item.english !== '' && (
                <Text style={styles.english}>{item.english}</Text>
              )}

              <View style={styles.metaRow}>
                <Text style={styles.restaurant}>{item.restaurant}</Text>
                <Text style={styles.distance}>{item.distance} m</Text>
              </View>
              {item.isDemo && (
                <Text style={{ color: '#777777', fontSize: 10, marginTop: 4 }}>
                  Fictional demo restaurant
                </Text>
              )}
            </View>
          </Pressable>
        )}
                ListEmptyComponent={
          <View style={styles.empty}>
            {loading ? (
              <Text style={styles.label}>Loading dishes…</Text>
            ) : error ? (
              <>
                <Text accessibilityRole="alert" style={styles.label}>
                  Could not load dishes
                </Text>
                <Text style={styles.summary}>{error}</Text>
                <Choice
                  label="Retry"
                  onPress={() => setRetryCount((value) => value + 1)}
                />
              </>
            ) : (
              <>
                <Text style={styles.label}>No dishes in this range</Text>
                <Text style={styles.summary}>
                  Try a higher budget or a wider radius.
                </Text>
              </>
            )}
          </View>
        }
        ListFooterComponent={
          <View style={styles.footer}>
            <View style={styles.pagination}>
              <Choice
                label="Previous"
                disabled={loading || page === 0}
                onPress={() => changePage(page - 1)}
              />
              <Text style={styles.muted}>Page {page + 1}</Text>
              <Choice
                label="Next"
                disabled={loading || error !== null || !hasNext}
                onPress={() => changePage(page + 1)}
              />
            </View>
            <Text style={styles.preview}>
              {usingDeviceLocation
                ? 'Distances from your last detected location'
                : previewArea ? 'Development preview · Distances from Tokyo Station' : `Distances from ${locationLabel}`}            </Text>
          </View>
        } />
    </View>
    {locationOpen && <LocationSheet
      area={{ ...searchLocation, label: locationLabel, device: usingDeviceLocation }}
      locating={locating}
      locationError={locationError}
      onLocate={useMyLocation}
      onClose={closeLocation}
      onApply={(area: SearchArea) => {
        setSearchLocation({ latitude: area.latitude, longitude: area.longitude });
        setLocationLabel(area.label);
        setUsingDeviceLocation(Boolean(area.device));
        setPreviewArea(false);
        changePage(0);
        closeLocation();
      }}
    />}
    {filtersOpen && (
    <FilterSheet
      budget={budget}
      radius={radius}
      onClose={() => setFiltersOpen(false)}
      onApply={(nextBudget, nextRadius) => {
        setBudget(nextBudget);
        setRadius(nextRadius);
        changePage(0);
        setFiltersOpen(false);
      }}
    />
  )}
  </SafeAreaView>;
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: '#FAFAFA'
  },

  container: {
    flex: 1,
    width: '100%',
    maxWidth: 580,
    alignSelf: 'center'
  },

  header: {
    paddingHorizontal: 16,
    paddingTop: 20,
    paddingBottom: 16,
    backgroundColor: '#FFFFFF',
    borderBottomWidth: 0.5,
    borderBottomColor: '#B8B8B8'
  },

  topRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 12
  },

  brand: {
    fontSize: 28,
    fontWeight: '800',
    letterSpacing: -1,
    color: '#151515',
    flexShrink: 1
  },

  button: {
    minHeight: 44,
    paddingHorizontal: 14,
    justifyContent: 'center',
    borderWidth: 0.5,
    borderColor: '#777777',
    borderRadius: 22
  },

  buttonText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#222222'
  },

  filterActive: {
    backgroundColor: '#EEEEEE'
  },

  summary: {
    fontSize: 12,
    color: '#666666',
    marginTop: 8,
    lineHeight: 18
  },

  sortRow: {
    marginTop: 18,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 8
  },

  selected: {
    backgroundColor: '#191919',
    borderColor: '#191919'
  },

  selectedText: {
    color: '#FFFFFF'
  },

  muted: {
    fontSize: 12,
    color: '#666666'
  },

  filterIconButton: {
    width: 44,
    height: 44,
    alignItems: 'center',
    justifyContent: 'center',
  },

  filterGlyph: {
    alignItems: 'center',
    gap: 5,
  },

  filterStripe: {
    height: 2,
    borderRadius: 1,
    backgroundColor: '#171717',
  },

  filterPanel: {
    marginTop: 16,
    padding: 14,
    backgroundColor: '#F7F7F7',
    borderWidth: 0.5,
    borderColor: '#B8B8B8',
    borderRadius: 14,
    gap: 12
  },

  label: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333333'
  },

  options: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8
  },

  list: {
    flex: 1
  },

  content: {
    padding: 14,
    paddingBottom: 32
  },

  columns: {
    gap: 12,
    marginBottom: 14
  },

  card: {
    flex: 1,
    maxWidth: '50%',
    backgroundColor: '#FFFFFF',
    borderRadius: 16,
    borderWidth: 0.5,
    borderColor: 'rgba(0,0,0,0.32)',
    overflow: 'hidden'
  },

  photo: {
    width: '100%',
    aspectRatio: 0.9,
    backgroundColor: '#EEEEEE',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 10,
    borderBottomWidth: 0.5,
    borderBottomColor: '#D4D4D4'
  },

  photoMark: {
    fontSize: 28,
    color: '#A0A0A0',
    fontWeight: '300'
  },

  photoLabel: {
    fontSize: 8,
    letterSpacing: 2,
    color: '#777777'
  },

  cardBody: {
    padding: 10,
    flex: 1

  },

  nameRow: {
    flexDirection: 'row',
    gap: 4,
    alignItems: 'flex-start'
  },

  japanese: {
    flex: 1,
    fontSize: 14,
    lineHeight: 21,
    fontWeight: '600',
    color: '#171717'
  },

  price: {
    fontSize: 16,
    lineHeight: 21,
    fontWeight: '700',
    color: '#171717'
  },

  english: {
    fontSize: 12,
    lineHeight: 18,
    color: '#666666',
    marginTop: 3
  },

  metaRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
    alignItems: 'flex-end',
    gap: 5,
    marginTop: 20
  },

  restaurant: {
    fontSize: 11,
    lineHeight: 16,
    color: '#666666',
    flexShrink: 1
  },

  distance: {
    fontSize: 11,
    lineHeight: 16,
    color: '#666666'

  },

  footer: {
    gap: 22,
    marginTop: 10

  },

  pagination: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    gap: 8
  },

  disabled: {
    opacity: 0.35

  },

  pressed: {
    opacity: 0.65

  },

  preview: {
    textAlign: 'center',
    fontSize: 11,
    color: '#777777',
    lineHeight: 18
  },

  empty: {
    paddingVertical: 50,
    alignItems: 'center'

  },

  locationSelector: {
    minHeight: 44,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  locationLink: {
    flexShrink: 1,
    fontSize: 13,
    fontWeight: '600',
    color: '#171717',
  },
  previewBadge: {
    fontSize: 10,
    color: '#707070',
    backgroundColor: '#F2F2F2',
    paddingHorizontal: 6,
    paddingVertical: 3,
    borderRadius: 5,
  },
  compactControls: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    alignItems: 'center',
    gap: 6,
  },
  filterSummary: {
    fontSize: 12,
    lineHeight: 18,
    color: '#666666',
    flexShrink: 1,
  },
  sortControl: {
    minHeight: 44,
    justifyContent: 'center',
    marginLeft: 'auto',
    paddingHorizontal: 4,
  },
  sortText: { fontSize: 12, color: '#666666' },
});

import { useEffect, useRef, useState } from 'react';
import { fetchNearbyDishes } from '@/api/dishes';
import { FlatList, Pressable, StyleSheet, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { FilterSheet } from '@/components/filter-sheet';

type Dish = { id: number; japanese: string; english: string; restaurant: string; price: number; distance: number };

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

  const [results, setResults] = useState<Dish[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hasNext, setHasNext] = useState(false);
  const [retryCount, setRetryCount] = useState(0);

  const list = useRef<FlatList<Dish>>(null);

  function changePage(next: number) {
    setPage(next);
    list.current?.scrollToOffset({ offset: 0, animated: true });
  }

  useEffect(() => {
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
            latitude: 35.6812,
            longitude: 139.7671,
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
  }, [sort, budget, radius, page, retryCount]);

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
        <View style={styles.sortRow}>
          <View style={styles.options}>
            <Choice label="Distance" selected={sort === 'distance'} onPress={() => { setSort('distance'); changePage(0); }} />
            <Choice label="Price" selected={sort === 'price'} onPress={() => { setSort('price'); changePage(0); }} />
          </View>
        </View>
       <Text style={styles.summary}>Within {radius} m · Up to {yen(budget)}</Text>
      </View>
      <FlatList ref={list} data={results} numColumns={2}
        keyExtractor={(dish) => String(dish.id)} style={styles.list} contentContainerStyle={styles.content} columnWrapperStyle={styles.columns}
        renderItem={({ item }) => <View style={styles.card}>
          <View style={styles.photo} accessibilityLabel="Dish photo placeholder">
            <Text style={styles.photoMark}>写真</Text><Text style={styles.photoLabel}>PHOTO PREVIEW</Text>
          </View>
          <View style={styles.cardBody}>
            <View style={styles.nameRow}><Text style={styles.japanese}>{item.japanese}</Text><Text style={styles.price}>{yen(item.price)}</Text></View>
            {item.english !== '' && (
              <Text style={styles.english}>{item.english}</Text>
            )}
            <View style={styles.metaRow}><Text style={styles.restaurant}>{item.restaurant}</Text><Text style={styles.distance}>{item.distance} m</Text></View>
          </View>
        </View>}
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
              Development preview · Distances from Tokyo Station
            </Text>
          </View>
        } />
    </View>
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
});

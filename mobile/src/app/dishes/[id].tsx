import { DishPhoto } from '@/components/dish-photo';
import { DishPriceHistory } from '@/components/dish-price-history';
import { getDemoDishImage } from '@/constants/demo-images';
import { useEffect, useState } from 'react';
import { router, useLocalSearchParams } from 'expo-router';
import {
  Linking,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import {
  fetchDish,
  fetchRestaurant,
  type Dish,
  type Restaurant,
} from '@/api/dishes';

type Details = {
  dish: Dish;
  restaurant: Restaurant;
};

export default function DishDetailsScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();

  const [details, setDetails] = useState<Details | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [mapError, setMapError] = useState<string | null>(null);
  const [retryCount, setRetryCount] = useState(0);

  useEffect(() => {
    setDetails(null);
    setError(null);
    setMapError(null);

    if (typeof id !== 'string' || !/^[1-9]\d*$/.test(id)) {
      setError('Invalid dish ID.');
      setLoading(false);
      return;
    }

    const controller = new AbortController();
    let active = true;
    let timedOut = false;

    setLoading(true);

    const timeout = setTimeout(() => {
      timedOut = true;
      controller.abort();
    }, 10000);

    async function loadDetails() {
      try {
        const dish = await fetchDish(id, controller.signal);
        const restaurant = await fetchRestaurant(
          dish.restaurant_id,
          controller.signal,
        );

        if (active) setDetails({ dish, restaurant });
      } catch (cause) {
        if (!active) return;

        setError(
          timedOut
            ? 'The request timed out. Please try again.'
            : cause instanceof Error
              ? cause.message
              : 'Could not load dish details.',
        );
      } finally {
        clearTimeout(timeout);
        if (active) setLoading(false);
      }
    }

    void loadDetails();

    return () => {
      active = false;
      clearTimeout(timeout);
      controller.abort();
    };
  }, [id, retryCount]);

  function goBack() {
    if (router.canGoBack()) {
      router.back();
    } else {
      router.replace('/');
    }
  }

  async function openMaps() {
    if (!details) return;

    setMapError(null);

    const { latitude, longitude } = details.restaurant;
    const query = encodeURIComponent(`${latitude},${longitude}`);

    try {
      await Linking.openURL(
        `https://www.google.com/maps/search/?api=1&query=${query}`,
      );
    } catch {
      setMapError('Could not open Maps. Please try again.');
    }
  }

  const dishName = details?.dish.name_ja ?? details?.dish.name;
  const restaurantName =
    details?.restaurant.name_ja ?? details?.restaurant.name;

    return (
    <SafeAreaView style={styles.screen}>
      <View style={styles.page}>
        <View style={styles.header}>
          <Pressable
            accessibilityRole="button"
            accessibilityLabel="Back to dishes"
            onPress={goBack}
            style={({ pressed }) => [
              styles.back,
              pressed && styles.pressed,
            ]}
          >
            <Text style={styles.backText}>‹ Back</Text>
          </Pressable>
        </View>

        <ScrollView contentContainerStyle={styles.content}>
          {loading ? (
            <View style={styles.status}>
              <Text style={styles.description}>Loading dish…</Text>
            </View>
          ) : error ? (
            <View style={styles.status}>
              <Text style={styles.title}>Could not load dish</Text>

              <Text accessibilityRole="alert" style={styles.description}>
                {error}
              </Text>

              <Pressable
                accessibilityRole="button"
                onPress={() => setRetryCount((value) => value + 1)}
                style={styles.button}
              >
                <Text style={styles.buttonText}>Retry</Text>
              </Pressable>
            </View>
          ) : details ? (
            <>
              <DishPhoto
                source={getDemoDishImage(
                  details.restaurant.name,
                  details.dish.name_en,
                )}
                label={details.dish.name_en ?? details.dish.name}
                aspectRatio={2}
              />

              <View style={styles.dishSection}>
                <View style={styles.nameRow}>
                  <View style={styles.names}>
                    <Text accessibilityRole="header" style={styles.title}>
                      {dishName}
                    </Text>

                    {details.dish.name_en &&
                      details.dish.name_en !== dishName && (
                        <Text style={styles.translation}>
                          {details.dish.name_en}
                        </Text>
                      )}
                  </View>

                  <Text style={styles.price}>
                    ¥{details.dish.price.toLocaleString('en-US')}
                  </Text>
                </View>
              </View>

              <View style={styles.lowerSection}>
                <Text accessibilityRole="header" style={styles.sectionLabel}>
                  Restaurant
                </Text>

                <View style={styles.card}>
                  <Text style={styles.restaurantName}>
                    {restaurantName}
                  </Text>

                  {details.restaurant.name_en &&
                    details.restaurant.name_en !== restaurantName && (
                      <Text style={styles.translation}>
                        {details.restaurant.name_en}
                      </Text>
                    )}
                  {details.restaurant.name.startsWith('[DEMO] ') && (
                    <Text
                      style={{
                        color: '#777777',
                        fontSize: 12,
                        lineHeight: 18,
                        marginBottom: 8,
                      }}
                    >
                      Fictional demo restaurant. The map pin is for testing only.
                    </Text>
                  )}

                  <Pressable
                    accessibilityRole="link"
                    accessibilityLabel={`${details.restaurant.address}. Open in Google Maps`}
                    onPress={openMaps}
                    style={({ pressed }) => [
                      styles.addressLink,
                      pressed && styles.pressed,
                    ]}
                  >
                    <Text style={styles.addressText}>
                      {details.restaurant.address} ↗
                    </Text>
                    <Text style={styles.addressHint}>Open in Google Maps</Text>
                  </Pressable>

                  {mapError && (
                    <Text
                      accessibilityRole="alert"
                      style={styles.description}
                    >
                      {mapError}
                    </Text>
                  )}
                </View>
                <DishPriceHistory
                  key={details.dish.id}
                  dishId={details.dish.id}
                />
              </View>
            </>
          ) : null}
        </ScrollView>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: '#FFFFFF',
  },
  page: {
    flex: 1,
    width: '100%',
    maxWidth: 580,
    alignSelf: 'center',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 16,
    paddingHorizontal: 16,
    minHeight: 64,
    backgroundColor: '#FFFFFF',
    borderBottomWidth: 0.5,
    borderBottomColor: '#D0D0D0',
  },
  back: {
    minHeight: 44,
    justifyContent: 'center',
    paddingRight: 8,
  },
  backText: {
    fontSize: 17,
    color: '#171717',
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#171717',
  },
  content: {
    flexGrow: 1,
    backgroundColor: '#F3F3F3',
  },
  status: {
    padding: 24,
    gap: 14,
  },
  photo: {
    width: '100%',
    aspectRatio: 2,
    backgroundColor: '#E8E8E8',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 10,
  },
  photoMark: {
    fontSize: 36,
    color: '#999999',
  },
  photoLabel: {
    fontSize: 12,
    color: '#777777',
  },
  dishSection: {
    paddingHorizontal: 24,
    paddingVertical: 26,
    backgroundColor: '#FFFFFF',
    borderBottomWidth: 0.5,
    borderBottomColor: '#D0D0D0',
  },
  nameRow: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: 16,
  },
  names: {
    flex: 1,
    gap: 6,
  },
  title: {
    fontSize: 26,
    lineHeight: 35,
    fontWeight: '700',
    color: '#171717',
  },
  translation: {
    fontSize: 15,
    lineHeight: 23,
    color: '#707070',
  },
  price: {
    fontSize: 30,
    lineHeight: 38,
    fontWeight: '700',
    color: '#171717',
    flexShrink: 1,
  },
  lowerSection: {
    padding: 20,
    paddingBottom: 32,
    gap: 14,
  },
  sectionLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: '#777777',
  },
  card: {
    padding: 20,
    gap: 8,
    backgroundColor: '#FFFFFF',
    borderRadius: 20,
    borderWidth: 0.5,
    borderColor: '#B8B8B8',
  },
  restaurantName: {
    fontSize: 20,
    lineHeight: 28,
    fontWeight: '600',
    color: '#171717',
  },
  description: {
    fontSize: 14,
    lineHeight: 22,
    color: '#707070',
  },
  button: {
    minHeight: 44,
    alignSelf: 'flex-start',
    justifyContent: 'center',
    paddingHorizontal: 16,
    borderRadius: 22,
    borderWidth: 0.5,
    borderColor: '#777777',
    marginTop: 8,
  },
  buttonText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#171717',
  },
  pressed: {
    opacity: 0.6,
  },
  addressLink: {
    alignSelf: 'stretch',
    minHeight: 44,
    justifyContent: 'center',
    paddingVertical: 8,
    gap: 4,
  },
  addressText: {
    fontSize: 14,
    lineHeight: 22,
    color: '#171717',
    textDecorationLine: 'underline',
  },
  addressHint: {
    fontSize: 12,
    lineHeight: 18,
    color: '#666666',
  },
});

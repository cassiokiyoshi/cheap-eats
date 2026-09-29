import { useEffect, useRef, useState } from 'react';
import { router } from 'expo-router';
import {
  KeyboardAvoidingView,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import {
  createDish,
  fetchRestaurants,
  type Restaurant,
} from '@/api/dishes';

export default function NewDishScreen() {
  const [restaurants, setRestaurants] = useState<Restaurant[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [retryCount, setRetryCount] = useState(0);

  const [query, setQuery] = useState('');
  const [restaurantId, setRestaurantId] = useState<number | null>(null);
  const [name, setName] = useState('');
  const [nameJA, setNameJA] = useState('');
  const [nameEN, setNameEN] = useState('');
  const [price, setPrice] = useState('');

  const [saving, setSaving] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const submitting = useRef(false);
  const mounted = useRef(false);
  const saveController = useRef<AbortController | null>(null);

  useEffect(() => {
    mounted.current = true;

    return () => {
      mounted.current = false;
      saveController.current?.abort();
    };
  }, []);

  useEffect(() => {
    const controller = new AbortController();
    let active = true;

    setLoading(true);
    setLoadError(null);

    const timer = setTimeout(() => controller.abort(), 10000);

    async function load() {
      try {
        const data = await fetchRestaurants(controller.signal);
        if (active) setRestaurants(data);
      } catch {
        if (active) {
          setLoadError('Could not load restaurants. Please try again.');
        }
      } finally {
        clearTimeout(timer);
        if (active) setLoading(false);
      }
    }

    void load();

    return () => {
      active = false;
      clearTimeout(timer);
      controller.abort();
    };
  }, [retryCount]);

  const selectedRestaurant = restaurants.find(
    (restaurant) => restaurant.id === restaurantId,
  );

  const search = query.trim().toLowerCase();

  const filteredRestaurants = restaurants.filter((restaurant) =>
    [
      restaurant.name,
      restaurant.name_ja,
      restaurant.name_en,
      restaurant.address,
    ]
      .filter(Boolean)
      .join(' ')
      .toLowerCase()
      .includes(search),
  );

  function goBack() {
    if (submitting.current) return;

    if (router.canGoBack()) {
      router.back();
    } else {
      router.replace('/');
    }
  }

  async function submit() {
    if (submitting.current) return;

    setSubmitError(null);

    if (!selectedRestaurant) {
      setSubmitError('Select a restaurant.');
      return;
    }

    if (!name.trim()) {
      setSubmitError('Enter the original dish name.');
      return;
    }

    const priceText = price.trim();
    const numericPrice = Number(priceText);

    if (
      !/^[0-9]+$/.test(priceText) ||
      !Number.isSafeInteger(numericPrice) ||
      numericPrice < 1 ||
      numericPrice > 2_147_483_647
    ) {
      setSubmitError(
        'Enter a whole-number price between ¥1 and ¥2,147,483,647, without commas.',
      );
      return;
    }

    submitting.current = true;
    setSaving(true);

    const controller = new AbortController();
    saveController.current = controller;

    const timer = setTimeout(() => controller.abort(), 15000);
    let saved = false;

    try {
      const dish = await createDish(
        {
          restaurant_id: selectedRestaurant.id,
          name,
          name_ja: nameJA,
          name_en: nameEN,
          price: numericPrice,
        },
        controller.signal,
      );

      saved = true;

      if (mounted.current) {
        router.replace({
          pathname: '/dishes/[id]',
          params: { id: String(dish.id) },
        });
      }
    } catch (cause) {
      if (mounted.current) {
        setSubmitError(
          saved
            ? 'The dish was saved, but its page could not open. Return to browsing instead of submitting again.'
            : cause instanceof Error
              ? cause.message
              : 'Could not confirm the save. Check for the dish before submitting again.',
        );
      }
    } finally {
      clearTimeout(timer);
      saveController.current = null;

      // Keep submission locked after a confirmed save.
      if (!saved) {
        submitting.current = false;
      }

      if (mounted.current) setSaving(false);
    }
  }

  return (
    <SafeAreaView style={styles.screen}>
      <KeyboardAvoidingView
        style={styles.page}
        behavior={Platform.OS === 'ios' ? 'padding' : undefined}
      >
        <ScrollView
          keyboardShouldPersistTaps="handled"
          contentContainerStyle={styles.content}
        >
          <Pressable
            accessibilityRole="button"
            disabled={saving}
            onPress={goBack}
            style={styles.textButton}
          >
            <Text style={styles.link}>‹ Back</Text>
          </Pressable>

          <Text accessibilityRole="header" style={styles.title}>
            Add a dish
          </Text>

          <Text style={styles.muted}>
            Development preview. Use a fictional demo restaurant for testing.
          </Text>

          <Text style={styles.label}>Restaurant *</Text>

          {loading ? (
            <Text style={styles.muted}>Loading restaurants…</Text>
          ) : loadError ? (
            <>
              <Text accessibilityRole="alert" style={styles.muted}>
                {loadError}
              </Text>
              <Pressable
                accessibilityRole="button"
                onPress={() => setRetryCount((value) => value + 1)}
                style={styles.textButton}
              >
                <Text style={styles.link}>Retry</Text>
              </Pressable>
            </>
          ) : selectedRestaurant ? (
            <View style={styles.restaurant}>
              <Text style={styles.label}>
                {selectedRestaurant.name_ja ?? selectedRestaurant.name}
              </Text>
              <Text style={styles.muted}>
                {selectedRestaurant.name_en}
              </Text>
              <Text style={styles.muted}>
                {selectedRestaurant.address}
              </Text>
              <Pressable
                accessibilityRole="button"
                disabled={saving}
                onPress={() => setRestaurantId(null)}
                style={styles.textButton}
              >
                <Text style={styles.link}>Change restaurant</Text>
              </Pressable>
            </View>
          ) : (
            <>
              <TextInput
                accessibilityLabel="Search restaurants"
                placeholder="Search by name or address"
                placeholderTextColor="#777777"
                value={query}
                onChangeText={setQuery}
                editable={!saving}
                style={styles.input}
              />

              {filteredRestaurants.slice(0, 10).map((restaurant) => (
                <Pressable
                  key={restaurant.id}
                  accessibilityRole="button"
                  disabled={saving}
                  onPress={() => setRestaurantId(restaurant.id)}
                  style={({ pressed }) => [
                    styles.restaurant,
                    pressed && styles.disabled,
                  ]}
                >
                  <Text style={styles.label}>
                    {restaurant.name_ja ?? restaurant.name}
                  </Text>
                  <Text style={styles.muted}>{restaurant.name_en}</Text>
                  <Text style={styles.muted}>{restaurant.address}</Text>
                  {restaurant.name.startsWith('[DEMO] ') && (
                    <Text style={styles.muted}>Fictional demo restaurant</Text>
                  )}
                </Pressable>
              ))}

              {filteredRestaurants.length === 0 && (
                <Text style={styles.muted}>
                  {restaurants.length === 0
                    ? 'No restaurants available yet.'
                    : 'No restaurants match your search.'}
                </Text>
              )}

              {filteredRestaurants.length > 10 && (
                <Text style={styles.muted}>
                  Showing 10 matches. Refine your search to find another restaurant.
                </Text>
              )}
            </>
          )}

          <Text style={styles.label}>Original dish name *</Text>
          <TextInput
            accessibilityLabel="Original dish name, required"
            placeholder="Name as written on the menu"
            placeholderTextColor="#777777"
            value={name}
            onChangeText={setName}
            editable={!saving}
            maxLength={200}
            style={styles.input}
          />

          <Text style={styles.label}>Japanese name</Text>
          <TextInput
            accessibilityLabel="Japanese name, optional"
            placeholder="例：醤油ラーメン"
            placeholderTextColor="#777777"
            value={nameJA}
            onChangeText={setNameJA}
            editable={!saving}
            maxLength={200}
            style={styles.input}
          />

          <Text style={styles.label}>English name</Text>
          <TextInput
            accessibilityLabel="English name, optional"
            placeholder="Example: Shoyu Ramen"
            placeholderTextColor="#777777"
            value={nameEN}
            onChangeText={setNameEN}
            editable={!saving}
            maxLength={200}
            style={styles.input}
          />

          <Text style={styles.muted}>
            Translations are optional. Automatic translation is not connected yet.
          </Text>

          <Text style={styles.label}>Price in JPY *</Text>
          <TextInput
            accessibilityLabel="Price in Japanese yen, required"
            placeholder="850"
            placeholderTextColor="#777777"
            value={price}
            onChangeText={setPrice}
            editable={!saving}
            keyboardType="number-pad"
            maxLength={10}
            style={styles.input}
          />

          {submitError && (
            <Text accessibilityRole="alert" style={styles.muted}>
              {submitError}
            </Text>
          )}

          <Pressable
            accessibilityRole="button"
            accessibilityState={{
              disabled: saving || loading || !selectedRestaurant,
              busy: saving,
            }}
            disabled={saving || loading || !selectedRestaurant}
            onPress={() => void submit()}
            style={({ pressed }) => [
              styles.submit,
              (saving || loading || !selectedRestaurant || pressed) &&
                styles.disabled,
            ]}
          >
            <Text style={styles.submitText}>
              {saving ? 'Saving…' : 'Add dish'}
            </Text>
          </Pressable>
        </ScrollView>
      </KeyboardAvoidingView>
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
  content: {
    padding: 24,
    paddingBottom: 48,
    gap: 12,
  },
  title: {
    fontSize: 28,
    fontWeight: '700',
    color: '#171717',
  },
  label: {
    fontSize: 14,
    fontWeight: '600',
    color: '#171717',
  },
  muted: {
    fontSize: 13,
    lineHeight: 20,
    color: '#707070',
  },
  input: {
    minHeight: 48,
    borderWidth: 0.5,
    borderColor: '#999999',
    borderRadius: 12,
    paddingHorizontal: 14,
    paddingVertical: 12,
    fontSize: 16,
    color: '#171717',
    backgroundColor: '#FFFFFF',
  },
  restaurant: {
    padding: 14,
    gap: 5,
    borderWidth: 0.5,
    borderColor: '#B8B8B8',
    borderRadius: 14,
    backgroundColor: '#FAFAFA',
  },
  textButton: {
    minHeight: 44,
    justifyContent: 'center',
    alignSelf: 'flex-start',
  },
  link: {
    fontSize: 14,
    color: '#171717',
    textDecorationLine: 'underline',
  },
  submit: {
    minHeight: 50,
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 14,
    backgroundColor: '#171717',
    marginTop: 8,
  },
  submitText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#FFFFFF',
  },
  disabled: {
    opacity: 0.5,
  },
});

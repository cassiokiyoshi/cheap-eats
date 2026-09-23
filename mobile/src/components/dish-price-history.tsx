import { useEffect, useState } from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';

import {
  fetchDishPriceHistory,
  type DishPriceChange,
} from '@/api/dishes';

const yen = (value: number) =>
  `¥${value.toLocaleString('en-US')}`;

const dateFormatter = new Intl.DateTimeFormat('en-GB', {
  day: 'numeric',
  month: 'short',
  year: 'numeric',
  timeZone: 'Asia/Tokyo',
});

export function DishPriceHistory({ dishId }: { dishId: number }) {
  const [history, setHistory] = useState<DishPriceChange[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [retryCount, setRetryCount] = useState(0);

  useEffect(() => {
    const controller = new AbortController();
    let active = true;
    let timedOut = false;

    setLoading(true);
    setError(null);
    setHistory([]);

    const timeout = setTimeout(() => {
      timedOut = true;
      controller.abort();
    }, 10000);

    async function load() {
      try {
        const changes = await fetchDishPriceHistory(
          dishId,
          controller.signal,
        );

        if (active) setHistory(changes);
      } catch (cause) {
        if (!active) return;

        setError(
          timedOut
            ? 'Price history took too long to load.'
            : cause instanceof Error
              ? cause.message
              : 'Could not load price history.',
        );
      } finally {
        clearTimeout(timeout);
        if (active) setLoading(false);
      }
    }

    void load();

    return () => {
      active = false;
      clearTimeout(timeout);
      controller.abort();
    };
  }, [dishId, retryCount]);

  return (
    <View style={styles.section}>
      <Text accessibilityRole="header" style={styles.heading}>
        Price history
      </Text>

      <View style={styles.card}>
        {loading ? (
          <Text style={styles.muted}>Loading price history…</Text>
        ) : error ? (
          <>
            <Text accessibilityRole="alert" style={styles.muted}>
              {error}
            </Text>

            <Pressable
              accessibilityRole="button"
              onPress={() => setRetryCount((value) => value + 1)}
              style={({ pressed }) => [
                styles.retry,
                pressed && styles.pressed,
              ]}
            >
              <Text style={styles.retryText}>Retry</Text>
            </Pressable>
          </>
        ) : history.length === 0 ? (
          <Text style={styles.muted}>
            No price changes recorded yet.
          </Text>
        ) : (
          <>
            <Text style={styles.muted}>
              Latest recorded changes · Dates in Tokyo time
            </Text>

            {history.map((change) => {
              const difference = change.new_price - change.old_price;
              const direction =
                difference > 0
                  ? 'Increase'
                  : difference < 0
                    ? 'Decrease'
                    : 'No change';

              return (
                <View key={change.id} style={styles.entry}>
                  <Text style={styles.date}>
                    {dateFormatter.format(new Date(change.changed_at))}
                  </Text>

                  <View style={styles.priceRow}>
                    <Text style={styles.previousPrice}>
                      {yen(change.old_price)}
                    </Text>
                    <Text style={styles.muted}>→</Text>
                    <Text style={styles.newPrice}>
                      {yen(change.new_price)}
                    </Text>
                  </View>

                  <Text style={styles.muted}>
                    {direction}
                    {difference !== 0
                      ? ` · ${yen(Math.abs(difference))}`
                      : ''}
                  </Text>
                </View>
              );
            })}

            {history.length === 20 && (
              <Text style={styles.muted}>
                Showing the latest 20 changes.
              </Text>
            )}
          </>
        )}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  section: {
    gap: 14,
    marginTop: 12,
  },
  heading: {
    fontSize: 14,
    fontWeight: '600',
    color: '#777777',
  },
  card: {
    padding: 20,
    gap: 14,
    backgroundColor: '#FFFFFF',
    borderRadius: 20,
    borderWidth: 0.5,
    borderColor: '#B8B8B8',
  },
  entry: {
    gap: 6,
    paddingTop: 14,
    borderTopWidth: 0.5,
    borderTopColor: '#DDDDDD',
  },
  date: {
    fontSize: 13,
    color: '#707070',
  },
  priceRow: {
    flexDirection: 'row',
    alignItems: 'center',
    flexWrap: 'wrap',
    gap: 10,
  },
  previousPrice: {
    fontSize: 17,
    color: '#707070',
  },
  newPrice: {
    fontSize: 19,
    fontWeight: '600',
    color: '#171717',
  },
  muted: {
    fontSize: 12,
    lineHeight: 19,
    color: '#707070',
  },
  retry: {
    alignSelf: 'flex-start',
    minHeight: 44,
    justifyContent: 'center',
    paddingHorizontal: 16,
    borderWidth: 0.5,
    borderColor: '#777777',
    borderRadius: 22,
  },
  retryText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#171717',
  },
  pressed: {
    opacity: 0.6,
  },
});

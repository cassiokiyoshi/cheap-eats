import { StyleSheet, Text, View } from 'react-native';
import Svg, { Circle, Line, Polyline } from 'react-native-svg';

import type { DishPriceChange } from '@/api/dishes';

type Props = {
  history: DishPriceChange[];
};

const WIDTH = 320;
const HEIGHT = 140;
const PADDING = 14;

const yen = (value: number) =>
  `¥${value.toLocaleString('en-US')}`;

export function PriceHistoryChart({ history }: Props) {
  if (history.length === 0) return null;

  // The API returns newest first; the chart reads oldest first.
  const chronological = [...history].reverse();
  const oldest = chronological[0];

  // Include the price immediately before the first displayed change.
  const prices = [
    oldest.old_price,
    ...chronological.map((change) => change.new_price),
  ];

  const minimum = Math.min(...prices);
  const maximum = Math.max(...prices);
  const range = maximum - minimum;

  const points = prices.map((price, index) => ({
    x:
      PADDING +
      (index / (prices.length - 1)) * (WIDTH - PADDING * 2),
    y:
      range === 0
        ? HEIGHT / 2
        : HEIGHT -
          PADDING -
          ((price - minimum) / range) * (HEIGHT - PADDING * 2),
  }));

  const firstPrice = prices[0];
  const lastPrice = prices[prices.length - 1];
  const difference = lastPrice - firstPrice;

  const summary =
    difference === 0
      ? 'No net change across these records'
      : `${difference > 0 ? 'Increase' : 'Decrease'} of ${yen(
          Math.abs(difference),
        )} across these records`;

  return (
    <View style={styles.container}>
      <Text style={styles.caption}>
        Price changes on this page
      </Text>

      <View style={styles.rangeRow}>
        <Text style={styles.caption}>Low {yen(minimum)}</Text>
        <Text style={styles.caption}>High {yen(maximum)}</Text>
      </View>

      <View
        accessible
        accessibilityRole="image"
        accessibilityLabel={
          `Price chart, oldest to newest. Prices: ` +
          `${prices.map(yen).join(', ')}. ${summary}.`
        }
        style={styles.plot}
      >
        <Svg
          width="100%"
          height={HEIGHT}
          viewBox={`0 0 ${WIDTH} ${HEIGHT}`}
          accessible={false}
        >
          {[PADDING, HEIGHT / 2, HEIGHT - PADDING].map((y) => (
            <Line
              key={y}
              x1={PADDING}
              x2={WIDTH - PADDING}
              y1={y}
              y2={y}
              stroke="#DDDDDD"
              strokeWidth={1}
              strokeDasharray="3 5"
            />
          ))}

          <Polyline
            points={points.map(({ x, y }) => `${x},${y}`).join(' ')}
            fill="none"
            stroke="#171717"
            strokeWidth={2}
            strokeLinejoin="round"
            strokeLinecap="round"
          />

          {points.map(({ x, y }, index) => (
            <Circle
              key={index}
              cx={x}
              cy={y}
              r={index === points.length - 1 ? 4.5 : 3}
              fill={index === points.length - 1 ? '#171717' : '#FFFFFF'}
              stroke="#171717"
              strokeWidth={1.5}
            />
          ))}
        </Svg>
      </View>

      <View style={styles.rangeRow}>
        <View style={styles.endpoint}>
          <Text style={styles.price}>{yen(firstPrice)}</Text>
          <Text style={styles.caption}>Before first change</Text>
        </View>

        <View style={[styles.endpoint, styles.right]}>
          <Text style={styles.price}>{yen(lastPrice)}</Text>
          <Text style={styles.caption}>After last change</Text>
        </View>
      </View>

      <Text style={styles.summary}>{summary}</Text>

      <Text style={styles.caption}>
        Oldest → newest. Points show recorded changes, not elapsed time.
        Vertical scale adjusts to this page.
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    gap: 10,
  },
  rangeRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: 12,
  },
  plot: {
    backgroundColor: '#F5F5F5',
    borderRadius: 14,
    overflow: 'hidden',
  },
  endpoint: {
    flex: 1,
    gap: 3,
  },
  right: {
    alignItems: 'flex-end',
  },
  price: {
    fontSize: 17,
    fontWeight: '600',
    color: '#171717',
  },
  caption: {
    fontSize: 11,
    lineHeight: 17,
    color: '#707070',
  },
  summary: {
    fontSize: 13,
    lineHeight: 20,
    fontWeight: '500',
    color: '#171717',
  },
});

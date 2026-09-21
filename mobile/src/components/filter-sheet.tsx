import { useCallback, useMemo, useRef, useState } from 'react';
import {
  Animated,
  Modal,
  PanResponder,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
  useWindowDimensions,
} from 'react-native';
import Slider from '@react-native-community/slider';
import { useSafeAreaInsets } from 'react-native-safe-area-context';

type Props = {
  budget: number;
  radius: number;
  onClose: () => void;
  onApply: (budget: number, radius: number) => void;
};

const yen = (value: number) => `¥${value.toLocaleString('en-US')}`;

export function FilterSheet({
  budget,
  radius,
  onClose,
  onApply,
}: Props) {
  const [draftBudget, setDraftBudget] = useState(budget);
  const [draftRadius, setDraftRadius] = useState(radius);
  const insets = useSafeAreaInsets();
  const { height } = useWindowDimensions();
  const translateY = useRef(new Animated.Value(height)).current;
  const backdropOpacity = useRef(new Animated.Value(0)).current;
  const closing = useRef(false);
  const opened = useRef(false);

  const show = useCallback(() => {
    if (opened.current || closing.current) return;
    opened.current = true;

    Animated.parallel([
      Animated.timing(translateY, {
        toValue: 0,
        duration: 280,
        useNativeDriver: Platform.OS !== 'web',
      }),
      Animated.timing(backdropOpacity, {
        toValue: 1,
        duration: 220,
        useNativeDriver: Platform.OS !== 'web',
      }),
    ]).start();
  }, [translateY, backdropOpacity]);

  const dismiss = useCallback((afterClose?: () => void) => {
    if (closing.current) return;
    closing.current = true;

    Animated.parallel([
      Animated.timing(translateY, {
        toValue: height,
        duration: 250,
        useNativeDriver: Platform.OS !== 'web',
      }),
      Animated.timing(backdropOpacity, {
        toValue: 0,
        duration: 250,
        useNativeDriver: Platform.OS !== 'web',
      }),
    ]).start(({ finished }) => {
      if (finished) {
        if (afterClose) {
          afterClose();
        } else {
          onClose();
        }
      } else {
        closing.current = false;
      }
    });
  }, [height, onClose, translateY, backdropOpacity]);

  const panResponder = useMemo(() => {
    let moved = false;

    function springBack() {
      Animated.spring(translateY, {
        toValue: 0,
        useNativeDriver: Platform.OS !== 'web',
      }).start();
    }

    return PanResponder.create({
      onStartShouldSetPanResponderCapture: () => !closing.current,

      onMoveShouldSetPanResponderCapture: () => !closing.current,

      onPanResponderGrant: () => {
        moved = false;
        translateY.stopAnimation();
      },

      onPanResponderMove: (_, gesture) => {
        if (closing.current) return;

        if (Math.abs(gesture.dx) > 6 || Math.abs(gesture.dy) > 6) {
          moved = true;
        }

        translateY.setValue(Math.max(0, gesture.dy));
      },

      onPanResponderRelease: (_, gesture) => {
        if (closing.current) return;

        const draggedDown =
          gesture.dy > 80 ||
          (gesture.dy > 20 && gesture.vy > 0.8);

        if (!moved || draggedDown) {
          dismiss();
        } else {
          springBack();
        }
      },

      onPanResponderTerminationRequest: () => false,

      onPanResponderTerminate: () => {
        if (!closing.current) springBack();
      },
    });
  }, [dismiss, translateY]);

  return (
    <Modal
      visible
      transparent
      animationType="none"
      onShow={show}
      onRequestClose={() => dismiss()}
    >
      <View style={styles.overlay}>
        <Animated.View
          style={[styles.backdrop, { opacity: backdropOpacity }]}
        >
          <Pressable
            style={{ flex: 1 }}
            accessibilityRole="button"
            accessibilityLabel="Close filters without applying"
            onPress={() => dismiss()}
          />
        </Animated.View>

        <Animated.View
          accessibilityViewIsModal
          onAccessibilityEscape={() => dismiss()}
          style={[
            styles.sheet,
            {
              paddingBottom: Math.max(insets.bottom, 20),
              transform: [{ translateY }],
            },
          ]}
        >
          <View
            {...panResponder.panHandlers}
            style={{
              ...(Platform.OS === 'web'
                ? { touchAction: 'none' as const }
                : {}),
            }}
          >
            <Pressable
              accessibilityRole="button"
              accessibilityLabel="Close filters"
              accessibilityHint="Tap or drag down to close without applying."
              onPress={() => dismiss()}
              style={styles.handleArea}
            >
              <View style={styles.handle} />
            </Pressable>
          </View>

          <ScrollView
            contentContainerStyle={styles.content}
            bounces={false}
          >
            <Text accessibilityRole="header" style={styles.title}>
              Filters
            </Text>

            <View style={styles.section}>
              <View style={styles.labelRow}>
                <Text style={styles.label}>Maximum price</Text>
                <Text style={styles.value}>{yen(draftBudget)}</Text>
              </View>

              <Slider
                style={styles.slider}
                accessibilityLabel="Maximum price in yen"
                minimumValue={100}
                maximumValue={2000}
                step={100}
                value={draftBudget}
                onValueChange={setDraftBudget}
                minimumTrackTintColor="#171717"
                maximumTrackTintColor="#D5D5D5"
                thumbTintColor="#171717"
              />

              <View style={styles.labelRow}>
                <Text style={styles.hint}>¥100</Text>
                <Text style={styles.hint}>¥2,000</Text>
              </View>
            </View>

            <View style={styles.section}>
              <View style={styles.labelRow}>
                <Text style={styles.label}>Search radius</Text>
                <Text style={styles.value}>{draftRadius} m</Text>
              </View>

              <Slider
                style={styles.slider}
                accessibilityLabel="Search radius in meters"
                minimumValue={100}
                maximumValue={1000}
                step={100}
                value={draftRadius}
                onValueChange={setDraftRadius}
                minimumTrackTintColor="#171717"
                maximumTrackTintColor="#D5D5D5"
                thumbTintColor="#171717"
              />

              <View style={styles.labelRow}>
                <Text style={styles.hint}>100 m</Text>
                <Text style={styles.hint}>1,000 m</Text>
              </View>
            </View>

            <View style={styles.actions}>
              <Pressable
                accessibilityRole="button"
                onPress={() => {
                  setDraftBudget(1000);
                  setDraftRadius(300);
                }}
                style={({ pressed }) => [
                  styles.action,
                  pressed && styles.pressed,
                ]}
              >
                <Text style={styles.resetText}>Reset</Text>
              </Pressable>

              <Pressable
                accessibilityRole="button"
                onPress={() =>
                  dismiss(() => onApply(draftBudget, draftRadius))
                }
                style={({ pressed }) => [
                  styles.action,
                  styles.apply,
                  pressed && styles.pressed,
                ]}
              >
                <Text style={styles.applyText}>Apply</Text>
              </Pressable>
            </View>
          </ScrollView>
        </Animated.View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    justifyContent: 'flex-end',
  },

  backdrop: {
    position: 'absolute',
    top: 0,
    right: 0,
    bottom: 0,
    left: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.35)',
  },
  sheet: {
    width: '100%',
    maxWidth: 580,
    maxHeight: '90%',
    alignSelf: 'center',
    backgroundColor: '#FFFFFF',
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    overflow: 'hidden',
  },
  handleArea: {
    height: 44,
    alignItems: 'center',
    justifyContent: 'center',
  },
  handle: {
    width: 36,
    height: 4,
    borderRadius: 2,
    backgroundColor: '#858585',
  },
  content: {
    paddingHorizontal: 24,
    paddingBottom: 4,
    gap: 28,
  },
  title: {
    fontSize: 23,
    fontWeight: '700',
    color: '#171717',
  },
  section: {
    gap: 4,
  },
  labelRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    gap: 12,
  },
  label: {
    fontSize: 15,
    color: '#333333',
  },
  value: {
    fontSize: 17,
    fontWeight: '600',
    color: '#171717',
  },
  slider: {
    width: '100%',
    height: 44,
  },
  hint: {
    fontSize: 12,
    color: '#777777',
  },
  actions: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: 12,
  },
  action: {
    minHeight: 44,
    paddingHorizontal: 22,
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 22,
  },
  apply: {
    backgroundColor: '#171717',
  },
  resetText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333333',
  },
  applyText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#FFFFFF',
  },
  pressed: {
    opacity: 0.65,
  },
});

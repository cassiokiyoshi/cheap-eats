import { useEffect, useState } from 'react';
import {
  Image,
  StyleSheet,
  Text,
  View,
  type ImageSourcePropType,
} from 'react-native';

type Props = {
  source?: ImageSourcePropType;
  label: string;
  aspectRatio?: number;
};

export function DishPhoto({
  source,
  label,
  aspectRatio = 0.9,
}: Props) {
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    setFailed(false);
  }, [source]);

  const showImage = source !== undefined && !failed;

  return (
    <View style={[styles.container, { aspectRatio }]}>
      {showImage ? (
        <>
          <Image
            source={source}
            style={styles.image}
            resizeMode="cover"
            accessible
            accessibilityLabel={`${label}. AI-generated illustrative image.`}
            onError={() => setFailed(true)}
          />

          <View pointerEvents="none" style={styles.badge}>
            <Text style={styles.badgeText}>
              AI-generated · Illustrative only
            </Text>
          </View>
        </>
      ) : (
        <View style={styles.placeholder}>
          <Text style={styles.placeholderText}>No photo yet</Text>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    width: '100%',
    overflow: 'hidden',
    backgroundColor: '#EEEEEE',
  },
  image: {
    width: '100%',
    height: '100%',
  },
  badge: {
    position: 'absolute',
    bottom: 8,
    left: 8,
    right: 8,
    alignSelf: 'flex-start',
    backgroundColor: 'rgba(255,255,255,0.94)',
    borderRadius: 5,
    paddingHorizontal: 6,
    paddingVertical: 4,
  },
  badgeText: {
    color: '#333333',
    fontSize: 10,
    lineHeight: 14,
  },
  placeholder: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  placeholderText: {
    color: '#777777',
    fontSize: 12,
  },
});

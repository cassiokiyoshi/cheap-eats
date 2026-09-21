import type { ImageSourcePropType } from 'react-native';

const images: Record<string, ImageSourcePropType> = {
  'Shoyu Ramen': require('../../assets/demo/shoyu-ramen.png'),
  'Japanese Curry': require('../../assets/demo/japanese-curry.png'),
  Gyudon: require('../../assets/demo/gyudon.png'),
  Takoyaki: require('../../assets/demo/takoyaki.png'),
  Gyoza: require('../../assets/demo/gyoza.png'),
  Udon: require('../../assets/demo/udon.png'),
  'Zaru Soba': require('../../assets/demo/zaru-soba.png'),
  'Oyakodon Set': require('../../assets/demo/oyakodon.png'),
  'Karaage Set': require('../../assets/demo/karaage-set.png'),
  'Grilled Mackerel Set':
    require('../../assets/demo/grilled-mackerel-set.png'),
};

export function getDemoDishImage(
  restaurantName: string,
  dishName: string | null | undefined,
): ImageSourcePropType | undefined {
  if (!restaurantName.startsWith('[DEMO] ') || !dishName) {
    return undefined;
  }

  return Object.prototype.hasOwnProperty.call(images, dishName)
    ? images[dishName]
    : undefined;
}

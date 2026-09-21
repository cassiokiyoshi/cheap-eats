export type Dish = {
  id: number;
  restaurant_id: number;
  name: string;
  name_ja: string | null;
  name_en: string | null;
  price: number;
  currency: string;
};

export type Restaurant = {
  id: number;
  name: string;
  name_ja: string | null;
  name_en: string | null;
  address: string;
  latitude: number;
  longitude: number;
};

export type NearbyDish = {
  dish: Dish;
  restaurant: Restaurant;
  distance_meters: number;
};

export type NearbyDishOptions = {
  latitude: number;
  longitude: number;
  radius: number;
  maxPrice: number;
  sort: 'distance' | 'price';
  limit: number;
  offset: number;
};

export async function fetchNearbyDishes(
  options: NearbyDishOptions,
  signal?: AbortSignal,
): Promise<NearbyDish[]> {
  const baseURL = process.env.EXPO_PUBLIC_API_URL?.trim();

  if (!baseURL) {
    throw new Error('EXPO_PUBLIC_API_URL is not configured.');
  }

  const params = new URLSearchParams({
    lat: String(options.latitude),
    lng: String(options.longitude),
    radius: String(options.radius),
    max_price: String(options.maxPrice),
    sort: options.sort,
    limit: String(options.limit),
    offset: String(options.offset),
  });

  const response = await fetch(
    `${baseURL.replace(/\/+$/, '')}/api/dishes/nearby?${params.toString()}`,
    {
      headers: { Accept: 'application/json' },
      signal,
    },
  );

  if (!response.ok) {
    throw new Error(
      `Could not load nearby dishes (HTTP ${response.status}).`,
    );
  }

  const data: unknown = await response.json();

  if (!Array.isArray(data)) {
    throw new Error('The API returned an unexpected response.');
  }

  return data as NearbyDish[];
}

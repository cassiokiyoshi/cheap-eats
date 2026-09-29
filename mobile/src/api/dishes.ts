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

async function fetchDetail<T>(
  path: string,
  label: string,
  signal?: AbortSignal,
): Promise<T> {
  const baseURL = process.env.EXPO_PUBLIC_API_URL?.trim();

  if (!baseURL) {
    throw new Error('EXPO_PUBLIC_API_URL is not configured.');
  }

  const response = await fetch(
    `${baseURL.replace(/\/+$/, '')}${path}`,
    {
      headers: { Accept: 'application/json' },
      signal,
    },
  );

  if (response.status === 404) {
    throw new Error(`${label} not found.`);
  }

  if (!response.ok) {
    throw new Error(
      `Could not load ${label.toLowerCase()} (HTTP ${response.status}).`,
    );
  }

  return response.json() as Promise<T>;
}

export function fetchDish(id: string, signal?: AbortSignal) {
  return fetchDetail<Dish>(
    `/api/dishes/${encodeURIComponent(id)}`,
    'Dish',
    signal,
  );
}

export function fetchRestaurant(id: number, signal?: AbortSignal) {
  return fetchDetail<Restaurant>(
    `/api/restaurants/${id}`,
    'Restaurant',
    signal,
  );
}

export type DishPriceChange = {
  id: number;
  dish_id: number;
  old_price: number;
  new_price: number;
  currency: string;
  changed_at: string;
};

export async function fetchDishPriceHistory(
  dishId: number,
  limit: number,
  offset: number,
  signal?: AbortSignal,
): Promise<DishPriceChange[]> {
  const params = new URLSearchParams({
    limit: String(limit),
    offset: String(offset),
  });

  const data = await fetchDetail<unknown>(
    `/api/dishes/${dishId}/price-history?${params.toString()}`,
    'Price history',
    signal,
  );

  if (!Array.isArray(data)) {
    throw new Error('The API returned an unexpected price history.');
  }

  return data as DishPriceChange[];
}

export type CreateDishInput = {
  restaurant_id: number;
  name: string;
  name_ja?: string;
  name_en?: string;
  price: number;
};

export async function fetchRestaurants(
  signal?: AbortSignal,
): Promise<Restaurant[]> {
  const data = await fetchDetail<unknown>(
    '/api/restaurants',
    'Restaurants',
    signal,
  );

  if (!Array.isArray(data)) {
    throw new Error('The API returned an unexpected restaurant list.');
  }

  return data as Restaurant[];
}

export async function createDish(
  input: CreateDishInput,
  signal?: AbortSignal,
): Promise<Dish> {
  const baseURL = process.env.EXPO_PUBLIC_API_URL?.trim();

  if (!baseURL) {
    throw new Error('EXPO_PUBLIC_API_URL is not configured.');
  }

  const name = input.name.trim();

  if (!name) {
    throw new Error('Enter a dish name.');
  }

  if (
    !Number.isSafeInteger(input.restaurant_id) ||
    input.restaurant_id < 1
  ) {
    throw new Error('Select a restaurant.');
  }

  if (
    !Number.isInteger(input.price) ||
    input.price < 1 ||
    input.price > 2_147_483_647
  ) {
    throw new Error('Price must be a whole number between ¥1 and ¥2,147,483,647.');
  }

  let response: Response;

  try {
    response = await fetch(
      `${baseURL.replace(/\/+$/, '')}/api/dishes`,
      {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          restaurant_id: input.restaurant_id,
          name,
          name_ja: input.name_ja?.trim() || undefined,
          name_en: input.name_en?.trim() || undefined,
          price: input.price,
          currency: 'JPY',
        }),
        signal,
      },
    );
  } catch {
    throw new Error(
      'Could not confirm whether the dish was saved. Check the restaurant’s dishes before submitting again.',
    );
  }

  if (!response.ok) {
    if (response.status === 400 || response.status === 413) {
      const message = await response.text();
      throw new Error(message.trim() || 'Check the dish details.');
    }

    throw new Error(
      `Could not confirm the save (HTTP ${response.status}). Check the restaurant’s dishes before submitting again.`,
    );
  }

  let data: unknown;

  try {
    data = await response.json();
  } catch {
    throw new Error(
      'The server accepted the dish, but its response could not be read. Check the restaurant’s dishes before submitting again.',
    );
  }

  if (
    typeof data !== 'object' ||
    data === null ||
    !('id' in data) ||
    typeof data.id !== 'number' ||
    !Number.isSafeInteger(data.id) ||
    data.id < 1
  ) {
    throw new Error(
      'The server accepted the dish, but returned an unexpected response. Check the restaurant’s dishes before submitting again.',
    );
  }

  return data as Dish;
}

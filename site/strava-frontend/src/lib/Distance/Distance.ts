export const DistanceToLocal = (distance: number): number => {
  return DistanceToMiles(distance)
}

export const DistanceToMiles = (distance: number): number => {
  return distance/1609.34
}

export const DistanceToKm = (distance: number): number => {
  return distance/1000
}

// TODO: Use local storage to decide
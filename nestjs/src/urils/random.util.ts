export function getRandomNumber(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

export function getRandomSourceOfEventDiscovery(): string {
  const options = ['Social media', 'Friends', 'Found myself'];
  const randomIndex = getRandomNumber(0, options.length - 1);
  return options[randomIndex];
}

export function getRandomDate(start: Date, end: Date): Date {
  const startTime = start.getTime();
  const endTime = end.getTime();
  const randomTime = getRandomNumber(startTime, endTime);
  return new Date(randomTime);
}

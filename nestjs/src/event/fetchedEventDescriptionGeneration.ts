export function generateEventDescription(event) {
  const type = `${event.classifications[0].segment.name} - ${event.classifications[0].genre.name} (${event.classifications[0].subGenre.name})`;
  const date = event.dates.start.localDate || 'Unknown Date';
  const time = event.dates.start.localTime || 'Unknown Time';
  const status = event.dates.status.code;
  const venue = event._embedded.venues[0];
  const venueName = venue.name;
  const venueAddress = venue.address.line1;
  const venueCity = venue.city.name;
  const venueState = venue.state.name;
  const venuePostalCode = venue.postalCode;
  const venueCountry = venue.country.name;

  return (
    `**Event Type:** ${type}<br />` +
    `**Date and Time:** ${date} at ${time}<br />` +
    `**Event Status:** ${status}<br />` +
    `**Venue:** ${venueName}, ${venueAddress}, ${venueCity}, ${venueState}, ${venuePostalCode}, ${venueCountry}`
  );
}

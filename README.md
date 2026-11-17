# Tour Das Hugel

Hello and welcome to my scrapped together Tour Das Hugel results tracker. This code is _mostly_ organized, and is _mostly_ stable. If you find any bugs or issues, please leave [make an issue](https://github.com/Emyrk/stravacollector/issues) to help me out.

# Not an Austin Cyclist? What the heck is this?

Tour das Hugel is an annual cycling event in Austin. It is not an official event, yet hundreds of cyclists gather for it. Of the hundreds, a smaller fraction will complete the course in it's entirety.

The course is designed to hit as many of the toughest climbs in Austin in a singular ride. Making this ride very challenging, not only by it's raw metrics alone, 112mi with ~12k ft, but with the constant steep ramps to endure.

https://dashugel.bike/route/das-hugel

For this reason, many participants will partake in a portion of the route.


## Ok, so what is this website?

Since this is not an official event, there is no timings. No results page. No medals, no ribbons, no trophies, not even a gathering at the "finish". 

This websites serves as a ribbon to finishers. A medal to the competitive. And a results list for the spectator.

The website polls Strava user's cycling activities, and uses Strava's segments feature to determine if any given activity classifies as a Hugel. If it does, only the climbs are included in the finisher's result time. This ignores things like waiting at stop lights, bombing descents, or other activities that would encourage unsafe riding practices.



## Admin: refresh a Strava athlete profile

After deployment, Steven (Strava athlete ID `2661162`) can force a profile refresh:

```http
POST /api/v1/refresh-athlete/67733138
```

While signed into dashugel.bike, run this in the browser's developer console:

```js
const response = await fetch('/api/v1/refresh-athlete/67733138', {
  method: 'POST',
  credentials: 'same-origin',
});
console.log(response.status, await response.json());
```

For command-line requests, send your dashugel session JWT, not a Strava access token:

```sh
curl --fail-with-body -X POST \
  -H "Strava-Auth-JWT: ${DASHUGEL_SESSION_JWT}" \
  https://dashugel.bike/api/v1/refresh-athlete/67733138
```

The endpoint uses the target athlete's stored Strava authorization, refreshes and
saves rotated OAuth credentials when necessary, and updates the profile from
Strava. A successful response contains the refreshed profile. Leaderboard caches
on the serving instance are invalidated immediately; other instances retain their
normal cache expiry. This does not re-download activities or change Strava data.

Missing authorization returns `404`. Invalid or revoked authorization requires the
athlete to reconnect and returns `502`, as do other Strava fetch failures. Strava
profile rate limits return `429`; database failures return `500`. Error responses
do not include OAuth credentials.

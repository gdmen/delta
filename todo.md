- Handler test data should be in a separate file & should be auto-updated by the tests

- Add "events" that contain:
  - measurements
  - images / videos / etc
  - notes / writeups

- Hardcoded API token for writes

- Add failure case unit tests

- Export data (or just save sql).

- Have consistent starting date across associated graphs

- Add table / historical view of individual measurements

- Add to dashboard:
> lifts graph (with pictures)

> body weight graph (with pictures)

- Add parsers for / get data from:
> fitnotes bodyweight

> http://judojournal.garymenezes.com/profile/gary

> notebook

> apple health (bodyweight)

-----
> tournaments

> meditation

> sleep

> salary


Bugs:
- Unit test API calls calculate e.g. maxes up to the current date => the 'correct' result is constantly outdated
- Fitocracy units are contained inline & aren't in CSV headers. So e.g. Pull Ups with no weight don't have units anywhere inline. If an unweighted Pull Up is ingested first, the MeasurementType will have no units and then future (weighted) Pull Ups will be stored with no units.

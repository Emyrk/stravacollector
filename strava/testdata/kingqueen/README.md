# King & Queen of Austin Leaderboard - Mock Data

This directory contains mock data for the **King & Queen of Austin** leaderboard feature, which ranks athletes based on their performance across the top 50 cycling segments in Austin, TX.

## Overview

The King & Queen of Austin leaderboard is a competitive ranking system that:
- Tracks the top 50 most popular/challenging segments in Austin
- Awards points based on placement (top 10 finishes)
- Maintains separate leaderboards for male (King) and female (Queen) athletes
- Shows each athlete's KOM/QOM count, top finishes, and total points

## Files

### 1. `austin_segments.json`
Contains the 50 featured Austin segments with full details:
- **Segment IDs**: 7654321 - 7654370
- **Popular locations**: 360 North, Mount Bonnell, Lady Bird Lake, Shoal Creek, Barton Hills, Zilker Park, Westlake Hills, etc.
- **Segment types**: Sprints, climbs, rollers, TTs, and flat sections
- **Metadata**: Distance, elevation gain, grade, location coordinates, effort counts

#### Notable Segments:
- **Mount Bonnell Climb** (7654322): Cat 1 climb, 5.8% avg grade, 12.3% max
- **Westlake Dr Rollers** (7654327): 5.6km with 57m gain
- **Pennybacker Bridge Vista** (7654336): Steep climb to iconic viewpoint
- **The Rollingwood Climb** (7654370): Hardest segment, 11.5% max grade
- **Lamar Blvd North** (7654366): Longest segment at 8.8km

### 2. `athlete_performances.json`
Mock athlete data and segment effort records:

#### Athletes (15 total):
**Male Athletes:**
- Alex Rodriguez (#1000001) - Current King leader
- Marcus Johnson (#1000003) - Strong climber
- David Martinez (#1000005) - Sprint specialist
- Ryan Williams (#1000007)
- Chris Brown (#1000009)
- Michael Anderson (#1000013)
- James Wilson (#1000011)
- Daniel Moore (#1000015)

**Female Athletes:**
- Sarah Chen (#1000002) - Current Queen leader
- Emily Thompson (#1000004) - All-around strong
- Jessica Lee (#1000006)
- Amanda Garcia (#1000008)
- Olivia Davis (#1000010)
- Sophia Miller (#1000012)
- Isabella Taylor (#1000014)

#### Segment Efforts:
Sample efforts for the first 10 segments with:
- Athlete ID
- Elapsed time (seconds)
- KOM/QOM rank

### 3. `leaderboard_calculated.json`
The calculated leaderboard results showing final rankings.

#### Scoring System:
- **10 points** for 1st place (KOM/QOM)
- **9 points** for 2nd place
- **8 points** for 3rd place
- ... continuing down to 1 point for 10th place

#### King Leaderboard (Top 8):
1. **Alex Rodriguez** - 425 points, 32 KOMs (48/50 segments)
2. **Marcus Johnson** - 398 points, 28 KOMs (46/50 segments)
3. **David Martinez** - 367 points, 24 KOMs (45/50 segments)
4. **Ryan Williams** - 342 points, 21 KOMs (43/50 segments)
5. **Chris Brown** - 315 points, 18 KOMs (42/50 segments)
6. **Michael Anderson** - 298 points, 16 KOMs (40/50 segments)
7. **James Wilson** - 276 points, 14 KOMs (39/50 segments)
8. **Daniel Moore** - 251 points, 12 KOMs (37/50 segments)

#### Queen Leaderboard (Top 7):
1. **Sarah Chen** - 412 points, 31 QOMs (47/50 segments)
2. **Emily Thompson** - 385 points, 27 QOMs (45/50 segments)
3. **Jessica Lee** - 356 points, 23 QOMs (44/50 segments)
4. **Amanda Garcia** - 329 points, 20 QOMs (42/50 segments)
5. **Olivia Davis** - 305 points, 17 QOMs (41/50 segments)
6. **Sophia Miller** - 287 points, 15 QOMs (39/50 segments)
7. **Isabella Taylor** - 264 points, 13 QOMs (38/50 segments)

#### Segment Leaders:
Each segment lists the current King and Queen (athlete holding KOM/QOM).

## Usage

This mock data can be used for:
1. **Frontend Development**: Build UI components for leaderboard display
2. **API Testing**: Test leaderboard endpoints
3. **Database Schema Design**: Plan tables for segments, efforts, and rankings
4. **Algorithm Testing**: Validate scoring and ranking calculations
5. **Demo/Presentation**: Show example leaderboard functionality

## Data Structure Notes

### Segment Model
```go
type Segment struct {
    ID                 int64
    Name               string
    ActivityType       string   // "Ride"
    Distance           float64  // meters
    AverageGrade       float64  // percentage
    MaximumGrade       float64  // percentage
    ElevationHigh      float64  // meters
    ElevationLow       float64  // meters
    StartLatlng        []float64
    EndLatlng          []float64
    ClimbCategory      int32    // 0-5 (5 is HC)
    City               string
    State              string
    Country            string
    TotalElevationGain float64  // meters
    TotalEffortCount   int32
    TotalAthleteCount  int32
    TotalStarCount     int32
}
```

### Athlete Model
```go
type Athlete struct {
    AthleteID  int64
    FirstName  string
    LastName   string
    Gender     string  // "M" or "F"
    ProfilePic string
}
```

### Effort Model
```go
type SegmentEffort struct {
    AthleteID   int64
    SegmentID   int64
    ElapsedTime float64  // seconds
    KomRank     int32    // 1-based rank
    MovingTime  float64
    StartDate   time.Time
}
```

### Leaderboard Entry Model
```go
type LeaderboardEntry struct {
    Rank              int32
    AthleteID         int64
    FirstName         string
    LastName          string
    TotalPoints       int32
    SegmentsCompleted int32
    Koms              int32  // or Qoms for female
    Top5Finishes      int32
    Top10Finishes     int32
    ProfilePic        string
}
```

## Next Steps

To implement this feature:

1. **Database Schema**: Create tables for:
   - `austin_segments` - Featured segment list
   - `segment_efforts` - Athlete performances
   - `leaderboard_cache` - Calculated rankings

2. **API Endpoints**:
   - `GET /api/v1/king-queen-austin` - Get full leaderboard
   - `GET /api/v1/king-queen-austin/king` - King leaderboard only
   - `GET /api/v1/king-queen-austin/queen` - Queen leaderboard only
   - `GET /api/v1/king-queen-austin/segments` - Segment list with leaders
   - `GET /api/v1/king-queen-austin/athlete/{id}` - Athlete's detailed performance

3. **Background Jobs**:
   - Daily recalculation of leaderboard
   - Real-time updates when new efforts are recorded
   - Cache invalidation strategy

4. **Frontend Components**:
   - Leaderboard table with sortable columns
   - Athlete detail cards
   - Segment map view
   - Progress charts and comparisons

## Example API Response

```json
{
  "king_leaderboard": [...],
  "queen_leaderboard": [...],
  "total_segments": 50,
  "last_updated": "2024-11-13T00:00:00Z",
  "scoring_method": "points_based",
  "description": "Top 50 Austin cycling segments"
}
```

## Credits

Mock data generated for the Strava Collector project.
All segment names and locations are based on real Austin, TX cycling routes.
Athlete names and performance data are fictional.

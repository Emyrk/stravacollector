# King & Queen of Austin - Frontend Component

## Overview
This is the frontend page for the King & Queen of Austin leaderboard feature. It displays rankings for the top 50 cycling segments in Austin, TX.

## Files
- **KingQueenAustin.tsx** - Main page component with leaderboard tables
- **mockData.ts** - Mock data for development/testing
- **index.ts** - Export file

## Features
- 📊 Dual leaderboards (King & Queen) with tabbed interface
- 🏆 Trophy icons and rankings (gold/silver/bronze for top 3)
- 📈 Statistics cards showing key metrics
- 👥 Athlete profiles with avatars
- 🎨 Responsive design with dark mode support
- 📱 Mobile-friendly layout

## Usage

### Access the Page
Navigate to `/king-queen-austin` in the application.

### Mock Data
The component uses mock data from `mockData.ts` for:
- King leaderboard (8 athletes)
- Queen leaderboard (7 athletes)
- Featured segments (10 shown)
- Stats and rankings

### API Integration (Future)
To connect to real API:

1. Create an API hook:
```typescript
// In src/api/kingqueen.ts
export const useKingQueenLeaderboard = () => {
  return useQuery(['kingQueen'], async () => {
    const response = await axios.get('/api/v1/king-queen-austin');
    return response.data;
  });
};
```

2. Update the component:
```typescript
// In KingQueenAustin.tsx
const { data, isLoading, error } = useKingQueenLeaderboard();
```

## Components

### Main Page (KingQueenAustin)
- Header with trophy icon and title
- 4 stat cards (segments, scoring, leaders)
- Tabbed interface for King/Queen leaderboards
- "How It Works" explanation card

### Leaderboard Table
- Rank badges (trophy for 1st, medals for 2nd/3rd)
- Athlete name and avatar
- Total points (bold, large)
- KOMs/QOMs with fire icon
- Top 5 and Top 10 finishes
- Segments completed (x/50)

## Styling
Uses Chakra UI components with:
- Responsive design (mobile, tablet, desktop)
- Color mode support (light/dark)
- Orange accent color scheme
- Hover effects on table rows

## Mock Data Structure

### LeaderboardEntry
```typescript
{
  rank: number;
  athleteId: number;
  firstName: string;
  lastName: string;
  totalPoints: number;
  segmentsCompleted: number;
  komsQoms: number;
  top5Finishes: number;
  top10Finishes: number;
  profilePic?: string;
}
```

### Current Mock Leaders
**King:** Alex Rodriguez (425 pts, 32 KOMs)  
**Queen:** Sarah Chen (412 pts, 31 QOMs)

## Screenshots

### Leaderboard Table
- Trophy icons for top 3 athletes
- Athlete avatars and names
- Points breakdown

### Stats Cards
- Featured Segments: 50
- Scoring Method: 10-1 pts
- Leader points displayed

## Testing
```bash
# Start dev server
npm start

# Navigate to
http://localhost:3000/king-queen-austin
```

## Future Enhancements
- [ ] Add segment list view
- [ ] Implement athlete detail pages
- [ ] Add search/filter functionality
- [ ] Show position changes (↑↓ arrows)
- [ ] Add charts and visualizations
- [ ] Real-time updates via WebSocket
- [ ] Export leaderboard to PDF/CSV
- [ ] Share on social media

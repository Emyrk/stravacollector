// Mock data for King & Queen of Austin Leaderboard

export interface LeaderboardEntry {
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

export interface LeaderboardStats {
  totalSegments: number;
  lastUpdated: string;
  scoringMethod: string;
}

export const mockLeaderboardStats: LeaderboardStats = {
  totalSegments: 50,
  lastUpdated: '2024-11-13T00:00:00Z',
  scoringMethod: 'Points awarded for top 10 finishes (10 pts for 1st, 9 for 2nd, etc.)',
};

export const mockKingLeaderboard: LeaderboardEntry[] = [
  {
    rank: 1,
    athleteId: 1000001,
    firstName: 'Alex',
    lastName: 'Rodriguez',
    totalPoints: 425,
    segmentsCompleted: 48,
    komsQoms: 32,
    top5Finishes: 45,
    top10Finishes: 48,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000001/12345/1/large.jpg',
  },
  {
    rank: 2,
    athleteId: 1000003,
    firstName: 'Marcus',
    lastName: 'Johnson',
    totalPoints: 398,
    segmentsCompleted: 46,
    komsQoms: 28,
    top5Finishes: 42,
    top10Finishes: 45,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000003/12347/1/large.jpg',
  },
  {
    rank: 3,
    athleteId: 1000005,
    firstName: 'David',
    lastName: 'Martinez',
    totalPoints: 367,
    segmentsCompleted: 45,
    komsQoms: 24,
    top5Finishes: 39,
    top10Finishes: 43,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000005/12349/1/large.jpg',
  },
  {
    rank: 4,
    athleteId: 1000007,
    firstName: 'Ryan',
    lastName: 'Williams',
    totalPoints: 342,
    segmentsCompleted: 43,
    komsQoms: 21,
    top5Finishes: 36,
    top10Finishes: 41,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000007/12351/1/large.jpg',
  },
  {
    rank: 5,
    athleteId: 1000009,
    firstName: 'Chris',
    lastName: 'Brown',
    totalPoints: 315,
    segmentsCompleted: 42,
    komsQoms: 18,
    top5Finishes: 33,
    top10Finishes: 39,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000009/12353/1/large.jpg',
  },
  {
    rank: 6,
    athleteId: 1000013,
    firstName: 'Michael',
    lastName: 'Anderson',
    totalPoints: 298,
    segmentsCompleted: 40,
    komsQoms: 16,
    top5Finishes: 31,
    top10Finishes: 37,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000013/12357/1/large.jpg',
  },
  {
    rank: 7,
    athleteId: 1000011,
    firstName: 'James',
    lastName: 'Wilson',
    totalPoints: 276,
    segmentsCompleted: 39,
    komsQoms: 14,
    top5Finishes: 28,
    top10Finishes: 35,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000011/12355/1/large.jpg',
  },
  {
    rank: 8,
    athleteId: 1000015,
    firstName: 'Daniel',
    lastName: 'Moore',
    totalPoints: 251,
    segmentsCompleted: 37,
    komsQoms: 12,
    top5Finishes: 25,
    top10Finishes: 33,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000015/12359/1/large.jpg',
  },
];

export const mockQueenLeaderboard: LeaderboardEntry[] = [
  {
    rank: 1,
    athleteId: 1000002,
    firstName: 'Sarah',
    lastName: 'Chen',
    totalPoints: 412,
    segmentsCompleted: 47,
    komsQoms: 31,
    top5Finishes: 44,
    top10Finishes: 47,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000002/12346/1/large.jpg',
  },
  {
    rank: 2,
    athleteId: 1000004,
    firstName: 'Emily',
    lastName: 'Thompson',
    totalPoints: 385,
    segmentsCompleted: 45,
    komsQoms: 27,
    top5Finishes: 41,
    top10Finishes: 44,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000004/12348/1/large.jpg',
  },
  {
    rank: 3,
    athleteId: 1000006,
    firstName: 'Jessica',
    lastName: 'Lee',
    totalPoints: 356,
    segmentsCompleted: 44,
    komsQoms: 23,
    top5Finishes: 38,
    top10Finishes: 42,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000006/12350/1/large.jpg',
  },
  {
    rank: 4,
    athleteId: 1000008,
    firstName: 'Amanda',
    lastName: 'Garcia',
    totalPoints: 329,
    segmentsCompleted: 42,
    komsQoms: 20,
    top5Finishes: 35,
    top10Finishes: 40,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000008/12352/1/large.jpg',
  },
  {
    rank: 5,
    athleteId: 1000010,
    firstName: 'Olivia',
    lastName: 'Davis',
    totalPoints: 305,
    segmentsCompleted: 41,
    komsQoms: 17,
    top5Finishes: 32,
    top10Finishes: 38,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000010/12354/1/large.jpg',
  },
  {
    rank: 6,
    athleteId: 1000012,
    firstName: 'Sophia',
    lastName: 'Miller',
    totalPoints: 287,
    segmentsCompleted: 39,
    komsQoms: 15,
    top5Finishes: 29,
    top10Finishes: 36,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000012/12356/1/large.jpg',
  },
  {
    rank: 7,
    athleteId: 1000014,
    firstName: 'Isabella',
    lastName: 'Taylor',
    totalPoints: 264,
    segmentsCompleted: 38,
    komsQoms: 13,
    top5Finishes: 26,
    top10Finishes: 34,
    profilePic: 'https://dgalywyr863hv.cloudfront.net/pictures/athletes/1000014/12358/1/large.jpg',
  },
];

export interface FeaturedSegment {
  segmentId: number;
  name: string;
  distance: number;
  averageGrade: number;
  kingId: number;
  kingName: string;
  queenId: number;
  queenName: string;
}

export const mockFeaturedSegments: FeaturedSegment[] = [
  {
    segmentId: 7654321,
    name: '360 North Sprint',
    distance: 1245.8,
    averageGrade: 2.1,
    kingId: 1000001,
    kingName: 'Alex Rodriguez',
    queenId: 1000002,
    queenName: 'Sarah Chen',
  },
  {
    segmentId: 7654322,
    name: 'Mount Bonnell Climb',
    distance: 892.3,
    averageGrade: 5.8,
    kingId: 1000003,
    kingName: 'Marcus Johnson',
    queenId: 1000002,
    queenName: 'Sarah Chen',
  },
  {
    segmentId: 7654323,
    name: 'Shoal Creek Sprint',
    distance: 2145.6,
    averageGrade: 0.8,
    kingId: 1000001,
    kingName: 'Alex Rodriguez',
    queenId: 1000002,
    queenName: 'Sarah Chen',
  },
  {
    segmentId: 7654324,
    name: 'Lady Bird Lake Loop North',
    distance: 4823.1,
    averageGrade: 0.3,
    kingId: 1000003,
    kingName: 'Marcus Johnson',
    queenId: 1000004,
    queenName: 'Emily Thompson',
  },
  {
    segmentId: 7654325,
    name: 'Barton Creek Greenbelt',
    distance: 3456.2,
    averageGrade: 1.5,
    kingId: 1000003,
    kingName: 'Marcus Johnson',
    queenId: 1000002,
    queenName: 'Sarah Chen',
  },
  {
    segmentId: 7654326,
    name: 'The Zilker Hill',
    distance: 678.4,
    averageGrade: 4.2,
    kingId: 1000001,
    kingName: 'Alex Rodriguez',
    queenId: 1000002,
    queenName: 'Sarah Chen',
  },
  {
    segmentId: 7654327,
    name: 'Westlake Dr Rollers',
    distance: 5678.9,
    averageGrade: 2.8,
    kingId: 1000001,
    kingName: 'Alex Rodriguez',
    queenId: 1000002,
    queenName: 'Sarah Chen',
  },
  {
    segmentId: 7654336,
    name: 'Pennybacker Bridge Vista',
    distance: 1234.9,
    averageGrade: 4.5,
    kingId: 1000003,
    kingName: 'Marcus Johnson',
    queenId: 1000004,
    queenName: 'Emily Thompson',
  },
  {
    segmentId: 7654366,
    name: 'Lamar Blvd North',
    distance: 8765.4,
    averageGrade: 1.5,
    kingId: 1000001,
    kingName: 'Alex Rodriguez',
    queenId: 1000002,
    queenName: 'Sarah Chen',
  },
  {
    segmentId: 7654370,
    name: 'The Rollingwood Climb',
    distance: 2456.7,
    averageGrade: 4.2,
    kingId: 1000001,
    kingName: 'Alex Rodriguez',
    queenId: 1000002,
    queenName: 'Sarah Chen',
  },
];

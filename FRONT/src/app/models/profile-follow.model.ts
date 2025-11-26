export interface Profile {
    id: string;
    // Gateway salje "userId", ne "user_id"
    userId: string;

    // Gateway salje "firstName", ne "first_name"
    firstName: string;
    lastName: string;
    role: string;

    // Gateway salje "profilePicture"
    profilePicture: string;
    biography: string;
    motto: string;
}

export interface Stats {
    userId: string;
    // OVO JE BIO PROBLEM:
    followersCount: number; // Bilo je followers_count
    followingCount: number; // Bilo je following_count
}

export interface RecommendationRaw {
    userId: string;
    mutualConnections: number;
}

export interface FullProfileDisplay extends Profile {
    stats?: Stats;
    isFollowing?: boolean;
    mutuals?: number;
}
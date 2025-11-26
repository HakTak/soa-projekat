export interface Profile {
    id: string;
    userId: string;
    firstName: string;
    lastName: string;
    role: string; // 'ADMIN' | 'TOURIST' | 'GUIDE'
    profilePicture: string;
    biography: string;
    motto: string;
    isBlocked: boolean;
}

export interface UpdateProfileRequest {
    firstName: string;
    lastName: string;
    profilePicture: string;
    biography: string;
    motto: string;
}

export interface ProfileStats {
    userId: string;
    followersCount: string | number; // API might return string "0" or number 0
    followingCount: string | number;
}
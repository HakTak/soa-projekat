export interface User {
    id: string;
    username: string;
    email: string;
    role: string;
    blocked: boolean;
}

// Gateway vraća objekat koji sadrži niz korisnika, ne direktno niz
export interface UserListResponse {
    users: User[];
}
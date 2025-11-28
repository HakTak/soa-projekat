import { Keypoint } from "./keypoint";

export enum TourStatus {
    DRAFT = 'DRAFT',
    PUBLISHED = 'PUBLISHED',
    ARCHIVED  = 'ARCHIVED',
}

export interface Tour {
    id?: string;
    userName: string;
    title: string;
    description: string;
    difficulty: string;
    tags: string;
    status: TourStatus;
    price: number;
    keypoints: Keypoint[];
    publisedAt: Date | null;
    archivedAt: Date | null;
}
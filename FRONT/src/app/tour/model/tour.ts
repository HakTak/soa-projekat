import { Keypoint } from "./keypoint";

export enum TourStatus{
    DRAFT = 0,
    PUBLISHED = 1,
    ARCHIVED  = 2,
}

export interface Tour {
    id?: string;
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
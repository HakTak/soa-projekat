import { Keypoint } from "./keypoint";
export interface Tour {
    id?: string;
    title: string;
    description: string;
    difficulty: string;
    tags: string;
    status: string;
    price: number;
    keypoints: Keypoint[];
}
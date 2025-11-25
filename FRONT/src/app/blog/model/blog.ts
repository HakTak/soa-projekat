import { Comment } from "./comment";
export interface Blog {
    id: string;
    title: string;
    authorName: string;
    content: string;
    tags: string; 
    createdAt: Date;
    imageUrls: string[];
    comments: Comment[];
    likeCount: number
}
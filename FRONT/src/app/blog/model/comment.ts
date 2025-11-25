export interface Comment {
    id: string;
    postId: string;
    authorName: string;
    content: string;
    createdAt: Date;
    updatedAt: Date | null;
}
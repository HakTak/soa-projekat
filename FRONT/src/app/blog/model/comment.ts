export interface Comment {
    id: string;
    postId: string;
    authorName: string;
    text: string;
    createdAt: Date;
    updatedAt: Date | null;
}
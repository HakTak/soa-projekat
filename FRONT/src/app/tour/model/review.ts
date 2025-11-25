export interface Review {
  id: string;
  rating: number;
  comment: string;
  tourId: string;
  userId: string; // Pretpostavljam da imamo ID korisnika
  username?: string; // Ako backend salje username, super. Ako ne, koristimo placeholder.
  createdAt?: Date; // Datum kreiranja
}
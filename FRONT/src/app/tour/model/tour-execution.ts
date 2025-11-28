export enum TourExecutionStatus {
  ACTIVE = 0,
  COMPLETED = 1,
  ABANDONED = 2
}

export interface CompletedKeypoint {
  keyPointId: string;
  completedAt: Date; // ISO date string
}

export interface TouristPosition {
  lat: number;
  lng: number;
}

export interface TourExecution {
  id?: string;                        // UUID
  tourId: string;
  touristId: string;
  status: TourExecutionStatus;       // ENUM
  startedAt: Date | null;                // ISO date string
  lastActivityAt: Date | null;          // ISO date string
  finishedAt: Date | null;               // ISO date string
  completedKeyPoints: CompletedKeypoint[];
  currentPosition: TouristPosition;
}
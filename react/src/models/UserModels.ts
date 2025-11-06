export type LoginDto = {
  email: string;
  password: string;
};

export type RegisterDto = {
  email: string;
  password: string;
  confirmPassword: string;
};

export type ErrorMessage = {
  message: string;
  description: string;
};

export type Role = {
  role: string;
  description: string | null;
};

export type Session = {
  id: string;
  ip: string;
  createdAt: string;
};

export type User = {
  id: string;
  email: string;
  role: Role[];
  session: Session[];
  createdAt: string;
  updatedAt: string;
};

export type Image = {
  id: string;
  link: string;
};

export type CreateEventDto = {
  title: string;
  description: string;
  latitude: number;
  longitude: number;
  images: Image[];
  eventDate: Date | null;
};

// export type SerializableCreateEventDto = Omit<CreateEventDto, "eventDate"> & {
//   eventDate: string; // Replaced Date with string for serialization
// };

export type Event = CreateEventDto & {
  id: string;
  organizer: string;
  createdAt: Date;
  updatedAt: Date;
};

export type EventList = {
  events: Event[];
  totalCount: number;
};

export type CreateParticipantDto = {
  fullName: string;
  email: string;
  dateOfBirth: Date;
  sourceOfEventDiscovery: string;
  eventId: string;
};

export type Participant = Omit<CreateParticipantDto, "eventId"> & {
  id: string;
};

export type RegisterParticipantDto = {
  fullName: string;
  email: string;
  dateOfBirth: Date | null;
  sourceOfEventDiscovery: string;
};

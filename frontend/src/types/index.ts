export type Member = {
  id: number;
  name: string;
};

export type Group = {
  id: number;
  name: string;
  invite_code: string;
  members: Member[];
};

export type PreferenceKind = "like" | "dislike";
export type PreferenceCategory = "genre" | "ingredient";

export type Preference = {
  kind: PreferenceKind;
  category: PreferenceCategory;
  value: string;
};

export type MessageRole = "member" | "ai";

export type Message = {
  id: number;
  role: MessageRole;
  member_id: number;
  member_name: string;
  text: string;
  created_at: string;
};

export type Candidate = {
  id: string;
  name: string;
  area: string;
  genres: string[];
  budget: number;
  matched_all: boolean;
  violation_count: number;
};

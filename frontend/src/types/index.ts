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

export type AreaEntry = {
  large_code: string;
  large_name: string;
  middle_code: string;
  middle_name: string;
  small_code: string;
  small_name: string;
};

export type BudgetOption = {
  code: string;
  name: string;
};

export type SearchOptions = {
  areas: AreaEntry[];
  budgets: BudgetOption[];
};

// 店舗検索の絞り込み(すべて任意)。エリアはHotPepperのエリアコード、Budgetは予算コード。
export type SearchCriteria = {
  large_area?: string;
  middle_area?: string;
  small_area?: string;
  budget?: string;
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

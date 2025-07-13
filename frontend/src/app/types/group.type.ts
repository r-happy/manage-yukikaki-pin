import { User } from './user.type';

export type Group = {
  group_id: string;
  group_name: string;
  group_description: string;
  group_created_by_id: string;
  group_created_by: User;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
};

export type GroupMember = {
  group_member_id: string;
  group_id: string;
  user_id: string;
  user: User;
  not_allowed: boolean;
  admin: boolean;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
};

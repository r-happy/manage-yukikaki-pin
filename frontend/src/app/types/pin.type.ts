import { Group } from './group.type';
import { User } from './user.type';

export type Pin = {
  pin_id: string;
  group_id: string;
  group: Group;
  pin_type_id: string;
  pin_type: PinType;
  pin_name: string;
  pin_type_description: string;
  pin_created_by_id: string;
  pin_created_by: User;
  latitude: number;
  longitude: number;
  not_allowed: boolean;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
};

export type PinType = {
  pin_type_id: string;
  group_id: string;
  group: Group;
  pin_type_name: string;
  pin_type_description: string;
  pin_type_created_by_id: string;
  pin_type_created_by: User;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
};

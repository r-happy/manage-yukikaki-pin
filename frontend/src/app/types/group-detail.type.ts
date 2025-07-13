import { Group, GroupMember } from './group.type';
import { Pin, PinType } from './pin.type';

export interface PinWithPinType {
  pin: Pin;
  pin_type: PinType;
}

export interface GroupDetail {
  group: Group;
  group_members: GroupMember[];
  pins: PinWithPinType[];
}

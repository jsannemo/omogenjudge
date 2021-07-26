import { Store } from "pullstate";

const userIdKey = "user";
const profileKey = "profile";

type authData = {
  userId: number | null;
  profile: Profile | null;
};

export const AuthStore = new Store<authData>({
  userId: getUserId(),
  profile: getProfile(),
});

export function isAuthenticated(): boolean {
  return getProfile() != null;
}

export type Profile = {
  username: string;
  full_name: string;
}

export function setProfile(profile: Profile): void {
  AuthStore.update(s => {
    s.profile = profile;
  });
  window.localStorage.setItem(profileKey, JSON.stringify(profile));
}

export function getProfile(): Profile | null {
  const data = window.localStorage.getItem(profileKey);
  if (data !== null) {
    return JSON.parse(data);
  } else {
    return null;
  }
}

export function setUserId(userId: number): void {
  AuthStore.update(s => {
    if (userId != s.userId) {
      s.profile = null;
    }
    s.userId = userId;
  });
  window.localStorage.removeItem(profileKey);
  if (userId == 0) {
    window.localStorage.removeItem(userIdKey);
  } else {
    window.localStorage.setItem(userIdKey, userId.toString());
  }
}

export function clearUserId(): void {
  setUserId(0);
}

export function getUserId(): number | null {
  const userId = window.localStorage.getItem(userIdKey);
  if (userId !== null) {
    return parseInt(userId, 10);
  }
  return null;
}

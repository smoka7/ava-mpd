import endpoints from "./endpoints";
import { fetchOrFail, sendCommand } from "./helpers";

export const settingsStore = defineStore("settings", {
  state: () => {
    return {
      settings: {} as SettingsResponse,
    };
  },
  actions: {
    async getSettings() {
      const response = await fetchOrFail<SettingsResponse>(endpoints.setting);
      this.settings = response;
    },
    async sendCommandToSetting(payload: SettingCommand) {
      sendCommand(endpoints.setting, payload.command, payload.data);
      this.getSettings();
    },
  },
});

type DatabaseStats = {
  albums: number;
  artists: number;
  db_playtime: number;
  db_update: number;
  playtime: number;
  songs: number;
  uptime: number;
};

type Output = {
  ID: number;
  Name: string;
  Enabled: boolean;
};

type SettingsResponse = {
  DatabaseStats: DatabaseStats;
  DownloadCoverArt: boolean;
  Outputs: Array<Output>;
  ReplayGain: string;
};

type SettingCommand = {
  command: string;
  data?: { Value: number };
};

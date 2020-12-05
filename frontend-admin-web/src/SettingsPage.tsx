import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import React, { useState } from "react";
import { LinkListItem } from "./AppDrawer";
import { Button } from "@rmwc/button";
import { LinearProgress } from "@rmwc/linear-progress";
import { Typography } from "@rmwc/typography";
import {
  CollapsibleList,
  List,
  ListItemGraphic,
  SimpleListItem,
} from "@rmwc/list";
import { Icon } from "@rmwc/icon";
import { ReactComponent as ZermeloIcon } from "./zermelo-clean.svg";
import { Drawer, DrawerContent } from "@rmwc/drawer";
import { Route, Switch as RouterSwitch } from "react-router-dom";
import "@rmwc/switch/styles";
import UserSettings from "./settings/UserSettings";
import NetworkSettings from "./settings/NetworkSettings";
import OrganizationSettings from "./settings/OrganizationSettings";
import ZermeloSettings from "./settings/ZermeloSettings";
import { SettingPageProps } from "./settings/useSetting";

interface SettingsStore {
  [key: string]: object | undefined;
}

const SettingsPage: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [areContentsModified, setAreContentsModified] = useState(false);
  const [settingsStore, setSettingsStore] = useState({} as SettingsStore);
  const [saveChanges, setSaveChanges] = useState(
    () => undefined as (() => void) | undefined
  );
  const store = {
    store: settingsStore,
    update: (store: SettingsStore) =>
      setSettingsStore({ ...settingsStore, ...store }),
  };

  const settingsProps: SettingPageProps = {
    setIsLoading: setIsLoading,
    setIsModified: setAreContentsModified,
    settingsStore: store,
    setSaveChanges: setSaveChanges,
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        width: "100%",
        height: "100%",
      }}
    >
      <div
        style={{
          display: "flex",
          marginLeft: 32,
          marginTop: 16,
          marginRight: 16,
          height: 40,
          justifyContent: "space-between",
        }}
      >
        <h1 style={{ marginTop: 0 }}>Instellingen</h1>
        <div>
          <Button
            icon="checkmark"
            raised
            style={{
              backgroundColor: areContentsModified ? "#4CAF50" : undefined,
            }}
            disabled={!areContentsModified}
            onClick={() => saveChanges && saveChanges()}
          >
            Opslaan
          </Button>
        </div>
      </div>

      <Theme use={"background"} wrap>
        <Elevation
          z={16}
          style={{
            flexGrow: 1,
            margin: 16,
            borderRadius: 4,
            height: "100%",
            overflow: "hidden",
          }}
        >
          <LinearProgress
            style={{
              height: isLoading ? 4 : 0,
              transition: "0.3s",
              transitionDelay: isLoading ? "0s" : "300ms",
            }}
          />

          <div
            style={{
              display: "flex",
              alignItems: "flex-start",
              height: "100%",
            }}
          >
            <Drawer style={{ height: "100%" }}>
              <DrawerContent>
                <List>
                  <LinkListItem to="/settings/account">
                    <ListItemGraphic icon="account_circle" />
                    Mijn account
                  </LinkListItem>

                  <CollapsibleList
                    defaultOpen
                    handle={
                      <SimpleListItem
                        text="Mijn school"
                        graphic={<Icon icon="school" />}
                        metaIcon={<Icon icon="chevron_right" />}
                      />
                    }
                  >
                    <LinkListItem to="/settings/organization/information">
                      <ListItemGraphic icon="info" />
                      Informatie
                    </LinkListItem>

                    <LinkListItem to="/settings/organization/networking">
                      <ListItemGraphic icon="settings_ethernet" />
                      Netwerken
                    </LinkListItem>

                    <CollapsibleList
                      defaultOpen
                      handle={
                        <SimpleListItem
                          text="Koppelingen"
                          graphic={<Icon icon="link" />}
                          metaIcon={<Icon icon="chevron_right" />}
                        />
                      }
                    >
                      <LinkListItem to="/settings/organization/integration/zermelo">
                        <ListItemGraphic icon={<ZermeloIcon />} />
                        Zermelo
                      </LinkListItem>
                    </CollapsibleList>
                  </CollapsibleList>
                </List>
              </DrawerContent>
            </Drawer>

            <div
              style={{
                display: "flex",
                width: "100%",
                height: "100%",
              }}
            >
              <div
                style={{
                  display: "flex",
                  margin: 16,
                  flexDirection: "column",
                  width: "100%",
                  alignItems: "flex-start",
                }}
              >
                <RouterSwitch>
                  <Route exact path="/settings">
                    <Typography use="body1">
                      Selecteer een instellingscategorie uit het menu links
                    </Typography>
                  </Route>

                  <Route exact path="/settings/account">
                    <UserSettings {...settingsProps} />
                  </Route>

                  <Route exact path="/settings/organization/information">
                    <OrganizationSettings {...settingsProps} />
                  </Route>

                  <Route exact path="/settings/organization/networking">
                    <NetworkSettings {...settingsProps} />
                  </Route>

                  <Route
                    exact
                    path="/settings/organization/integration/zermelo"
                  >
                    <ZermeloSettings {...settingsProps} />
                  </Route>
                </RouterSwitch>
              </div>
            </div>
          </div>
        </Elevation>
      </Theme>
    </div>
  );
};

export default SettingsPage;

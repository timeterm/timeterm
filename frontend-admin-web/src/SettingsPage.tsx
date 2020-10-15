import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import React, { useEffect, useMemo, useState } from "react";
import { useQuery } from "react-query";
import DevicesPage, { fetchAuthnd } from "./DevicesPage";
import { snackbarQueue } from "./snackbarQueue";
import { UserResponse, LinkListItem } from "./AppDrawer";
import { Button } from "@rmwc/button";
import { TextField } from "@rmwc/textfield";
import { LinearProgress } from "@rmwc/linear-progress";
import { Typography } from "@rmwc/typography";
import {
  CollapsibleList,
  List,
  ListItem,
  ListItemGraphic,
  SimpleListItem,
} from "@rmwc/list";
import { Icon } from "@rmwc/icon";
import { ReactComponent as ZermeloIcon } from "./zermelo-clean.svg";
import { Drawer, DrawerContent, DrawerHeader } from "@rmwc/drawer";
import { Redirect, Route, Switch as RouterSwitch } from "react-router-dom";
import { Select } from "@rmwc/select";
import { Switch } from "@rmwc/switch";
import "@rmwc/switch/styles";
import LoginDonePage from "./LoginDonePage";
import LoginPage from "./LoginPage";
import StudentsPage from "./StudentsPage";

interface OrganizationResponse {
  id: string;
  name: string;
}

const SettingsPage: React.FC = () => {
  const [organizationName, setOrganizationName] = useState("");

  const { isLoading, error, data: organization } = useQuery<
    OrganizationResponse
  >(
    "organizationInfo",
    () =>
      fetchAuthnd("/api/user/me")
        .then((res) => res.json())
        .then((user: UserResponse) =>
          fetchAuthnd(`/api/organization/${user.organizationId}`).then((res) =>
            res.json()
          )
        ),
    {
      refetchInterval: false,
      refetchOnMount: false,
      refetchOnReconnect: false,
      refetchOnWindowFocus: false,
    }
  );

  useEffect(() => {
    if (organization) setOrganizationName(organization?.name);
  }, [organization]);

  const areContentsModified = useMemo(() => {
    return organization && organizationName !== organization?.name;
  }, [organization, organizationName]);

  useEffect(() => {
    if (error)
      snackbarQueue.notify({
        title: <b>Er is een fout opgetreden</b>,
        body: "Kon data niet van server ophalen",
        icon: "error",
        dismissesOnAction: true,
        actions: [
          {
            title: "Sluiten",
            icon: "close",
          },
        ],
      });
  }, [error]);

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

                    <LinkListItem to="/settings/organization/wifi">
                      <ListItemGraphic icon="wifi" />
                      Wi-Fi
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
                margin: 16,
                flexDirection: "column",
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
                  <Typography use="headline5">Mijn account</Typography>

                  <TextField
                    style={{
                      marginTop: 16,
                      width: "25em",
                    }}
                    label={"Naam"}
                    outlined
                  />
                </Route>

                <Route exact path="/settings/organization/information">
                  <Typography use="headline5">Informatie</Typography>

                  <span style={{ marginTop: 8 }}>
                    Je bent op het moment onderdeel van de organisatie (school)
                    &quot;
                    {organization?.name}&quot;
                  </span>

                  <TextField
                    label={"Naam"}
                    value={organizationName}
                    style={{
                      marginTop: 16,
                    }}
                    outlined
                    onChange={(evt) => {
                      setOrganizationName(
                        (evt.target as HTMLInputElement).value
                      );
                    }}
                  />
                </Route>

                <Route exact path="/settings/organization/wifi">
                  <div
                    style={{
                      display: "flex",
                      flexDirection: "column",
                    }}
                  >
                    <Typography use="headline5">Wifi-instellingen</Typography>

                    <TextField
                      style={{
                        marginTop: 16,
                      }}
                      label={"Netwerknaam (SSID)"}
                      outlined
                    />

                    <div
                      style={{
                        marginTop: 16,
                      }}
                    >
                      <Select
                        label={"EAP-methode"}
                        enhanced
                        outlined
                        options={["PEAP"]}
                      />
                    </div>

                    <div
                      style={{
                        marginTop: 16,
                      }}
                    >
                      <Select
                        label={"Phase 2-verificatie"}
                        enhanced
                        outlined
                        options={["Geen", "MSCHAPV2", "GTC", "AKA", "AKA'"]}
                      />
                    </div>

                    <TextField
                      style={{
                        marginTop: 16,
                      }}
                      label={"Identiteit"}
                      outlined
                    />

                    <TextField
                      style={{
                        marginTop: 16,
                      }}
                      label={"Anonieme identiteit"}
                      outlined
                    />

                    <TextField
                      style={{
                        marginTop: 16,
                      }}
                      label={"Wachtwoord"}
                      outlined
                    />

                    <Switch
                      style={{
                        marginTop: 16,
                      }}
                    >
                      Netwerk is verborgen
                    </Switch>
                  </div>
                </Route>

                <Route exact path="/settings/organization/integration/zermelo">
                  <div
                    style={{
                      display: "flex",
                      flexDirection: "column",
                    }}
                  >
                    <Typography use="headline5">Zermelo-koppeling</Typography>

                    <TextField
                      style={{
                        marginTop: 16,
                        width: "25em",
                      }}
                      label={"Zermelo-institutie"}
                      outlined
                    />

                    <TextField
                      style={{
                        marginTop: 16,
                        width: "25em",
                      }}
                      label={"Token van Timeterm-gebruiker in Zermelo"}
                      outlined
                    />
                  </div>
                </Route>
              </RouterSwitch>
            </div>
          </div>
        </Elevation>
      </Theme>
    </div>
  );
};

export default SettingsPage;

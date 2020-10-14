import { Drawer, DrawerContent, DrawerHeader } from "@rmwc/drawer";
import Logo from "./logo-white.svg";
import { Theme, ThemeProvider } from "@rmwc/theme";
import {
  List,
  ListItem,
  ListItemGraphic,
  ListItemPrimaryText,
  ListItemSecondaryText,
  ListItemText,
} from "@rmwc/list";
import React, { useEffect } from "react";
import { Link, useHistory, useLocation } from "react-router-dom";
import Cookies from "universal-cookie";
import { useQuery } from "react-query";
import { fetchAuthnd } from "./DevicesPage";
import { snackbarQueue } from "./snackbarQueue";
import "@rmwc/snackbar/styles";

interface LinkListItemProps {
  to: string;
  className?: string;
}

const LinkListItem: React.FC<LinkListItemProps> = (props) => {
  const location = useLocation();

  return (
    <ListItem
      selected={location.pathname === props.to}
      className={props.className}
    >
      <Link
        to={props.to}
        style={{
          color: "inherit",
          textDecoration: "inherit",
          display: "flex",
          height: "100%",
          width: "100%",
          alignItems: "center",
        }}
      >
        {props.children}
      </Link>
    </ListItem>
  );
};

interface AnchorListItemProps
  extends React.AnchorHTMLAttributes<HTMLAnchorElement> {}

const AnchorListItem: React.FC<AnchorListItemProps> = (props) => {
  return (
    <ListItem className={props.className}>
      <a
        {...props}
        style={{
          color: "inherit",
          textDecoration: "inherit",
          display: "flex",
          height: "100%",
          width: "100%",
          alignItems: "center",
        }}
      >
        {props.children}
      </a>
    </ListItem>
  );
};

interface UserResponse {
  id: string;
  name: string;
  email: string;
}

const AppDrawer: React.FC = () => {
  const { isLoading, error, data: user } = useQuery<UserResponse>(
    "userInfo",
    () => fetchAuthnd("/api/user/me").then((res) => res.json())
  );
  const history = useHistory();

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
    <Theme use={["primaryBg", "onPrimary"]} wrap>
      <Drawer>
        <DrawerHeader
          style={{
            marginTop: 16,
          }}
        >
          <img src={Logo} alt={"Timeterm Logo"} width={96} />
        </DrawerHeader>
        <DrawerContent>
          <ThemeProvider
            options={{
              primary: "rgb(0, 0, 0)",
            }}
            style={{
              display: "flex",
              flexDirection: "column",
              height: "100%",
            }}
          >
            <List style={{ flexGrow: 1 }}>
              <Theme use={["onPrimary"]} wrap>
                <LinkListItem to={"/devices"}>
                  <Theme use={["onPrimary"]} wrap>
                    <ListItemGraphic icon="tablet" />
                  </Theme>
                  Apparaten
                </LinkListItem>
              </Theme>
              <Theme use={["onPrimary"]} wrap>
                <LinkListItem to={"/students"}>
                  <Theme use={["onPrimary"]} wrap>
                    <ListItemGraphic icon="group" />
                  </Theme>
                  Leerlingen
                </LinkListItem>
              </Theme>
              <Theme use={["onPrimary"]} wrap>
                <AnchorListItem
                  href={`timeterm:${btoa(new Cookies().get("ttsess"))}`}
                  rel={"noopener noreferrer"}
                >
                  <Theme use={["onPrimary"]} wrap>
                    <ListItemGraphic icon="bluetooth_connected" />
                  </Theme>
                  Apparaat koppelen
                </AnchorListItem>
              </Theme>
            </List>

            <List twoLine={true} style={{ paddingBottom: 0 }}>
              <Theme use={["onPrimary"]} wrap>
                <ListItem>
                  <Theme use={["onPrimary"]} wrap>
                    <ListItemGraphic icon="person" />
                  </Theme>
                  <ListItemText>
                    <ListItemPrimaryText>
                      {!isLoading && user?.name}
                    </ListItemPrimaryText>
                    <Theme use={["onPrimary"]} wrap>
                      <ListItemSecondaryText>
                        {!isLoading && user?.email}
                      </ListItemSecondaryText>
                    </Theme>
                  </ListItemText>
                </ListItem>
              </Theme>
            </List>
            <List>
              <Theme use={["onPrimary"]} wrap>
                <ListItem
                  onClick={() => {
                    new Cookies().remove("ttsess", {
                      path: "/",
                    });
                    history.push("/");
                  }}
                >
                  <Theme use={["onPrimary"]} wrap>
                    <ListItemGraphic icon="logout" />
                  </Theme>
                  Uitloggen
                </ListItem>
              </Theme>
            </List>
          </ThemeProvider>
        </DrawerContent>
      </Drawer>
    </Theme>
  );
};

export default AppDrawer;

import { History, Home, RefreshCw, Sparkles } from "lucide-react";
import SidebarSearch from "./SidebarSearch";
import SidebarItem from "./SidebarItem";
import SidebarSettings from "./SidebarSettings";

export default function Sidebar() {
  const iconSize = 20;

  return (
    <div
      className="
            z-50
            flex
            sm:flex-col
            justify-between
            sm:fixed
            sm:top-0
            sm:left-0
            sm:h-screen
            h-auto
            sm:w-auto
            w-full
            border-b
            sm:border-b-0
            sm:border-r
            border-(--color-bg-tertiary)
            pt-2
            sm:py-10
            sm:px-1
            px-4
            bg-(--color-bg)
        "
    >
      <div className="flex gap-4 sm:flex-col">
        <SidebarItem
          space={10}
          to="/"
          name="Home"
          onClick={() => {}}
          modal={<></>}
        >
          <Home size={iconSize} />
        </SidebarItem>
        <SidebarSearch size={iconSize} />
        <SidebarItem
          space={10}
          to="/recommendations"
          name="Recommendations"
          onClick={() => {}}
          modal={<></>}
        >
          <RefreshCw size={iconSize} />
        </SidebarItem>
        <SidebarItem
          space={10}
          to="/wrapped"
          name="Wrapped"
          onClick={() => {}}
          modal={<></>}
        >
          <Sparkles size={iconSize} />
        </SidebarItem>
        <SidebarItem
          space={10}
          to="/rewind"
          name="Rewind"
          onClick={() => {}}
          modal={<></>}
        >
          <History size={iconSize} />
        </SidebarItem>
      </div>
      <div className="flex gap-4 sm:flex-col">
        <SidebarSettings size={iconSize} />
      </div>
    </div>
  );
}

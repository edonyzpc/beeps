import clsx from "clsx";
import ClustrBar from "../ClustrBar";
import TagsSection from "../HomeSidebar/TagsSection";
import SearchBar from "../SearchBar";

interface Props {
  className?: string;
}

const TimelineSidebar = (props: Props) => {
  return (
    <aside
      className={clsx(
        "relative w-full h-auto max-h-screen overflow-auto hide-scrollbar flex flex-col justify-start items-start",
        props.className,
      )}
    >
      <SearchBar />
      <ClustrBar user={"timeline"} />
      <TagsSection />
    </aside>
  );
};

export default TimelineSidebar;

import { Tokens } from "marked";
import { Context } from "../parser/index.js";
import { Fetcher } from "../parser/fetcher.js";

// Custom events 
//
// Reference
// - https://stackoverflow.com/a/68783088
// - https://github.com/microsoft/TypeScript/issues/28357#issuecomment-748550734


declare global {
  interface GlobalEventHandlersEventMap {
    "state-changed": CustomEvent<{
      state: State
      push: boolean | null
    }>;
    "filter-click": import("../static/components/rpr-filter.js").FilterClickEvent;
    "result-click": import("../static/components/app-root.js").ResultClickEvent;
    "search-change": import("../static/components/app-root.js").SearchChangeEvent;
    "search-buffer": import("../static/components/app-root.js").SearchBufferEvent;
  }

  type SearchBufferAction = "open-buffer" | "reset-and-open-buffer" | "flush-and-close-buffer";

  type Filter = {
    key: string,
    value: string,
  }

  type State = {
    search: string;
    filters: Filter[];
    path: string;
    hash: string;
  }

  type Issue = {
    id: number;
    title: string;
  }

  type Section = {
    title: string;
    depth: number;
    id: string;
    prefix: string;
    timecode: string;
    issues: Issue[];
    excerpt: string | null;
    subsections: Section[];
  }

  type Contributor = {
    name: string;
    filterable: boolean;
    initials: string;
    avatar: string;
  }

  /**
   * Data for a visual search result.
   *
   * This describes a markdown asset, but also includes search meta data.
   */
  type Result = {
    /** The asset title from the markdown frontmatter */
    title: string;

    /** The asset slug from the markdown file name */
    slug: string;

    /** The asset date from the markdonw file name */
    date: string;

    /** The name of the podcast, person, etc. that published this asset */
    publisher: string;

    /** The Pagefind result ID corresponding to this Result */
    pagefindResultId: string;

    /** True if the current Pagefind data has been loaded for this result */
    loaded: Boolean;

    /** 
    * Excerpt pertenant to the current search string that occurs before the 
    * first section heading. 
    *
    * Each document section in [sections] has it's own excerpt.
    */
    excerpt: string | null;

    /**
     * Issues that are references before the first section heading.
     *
     * Issues are visual callouts in the rendered HTML of an asset that 
     * link to a GitHub issue tracking a potential improvement to the text.
     *
     * Each section in [sections] has it's own issues, these issues are only
     * the ones defined before the first section heading.
     */
    issues: Issue[];

    /**
     * Each markdown heading becomes an "section" of the asset. 
     *
     * For example, if the asset medium is an "article" each heading in the 
     * asset becomes a section which can be displayed in search results.
     *
     * For "audio" assets, headings mark changes in the topic of conversation,
     * and have an associated timecode.
     */
    sections: Section[];

    /**
     * The relevance of this result to the current search (higher = more 
     * relevant).
     *
     * An undefined score means the result did not match the search criteria.
     */
    score: number | undefined;

    /**
     * The visual order of this result within the list of all results. Must be 
     * an integer, smallest first.
     */
    order: number | undefined;
  }

  type FrontMatterSource = {
    url?: string;
    kind: string;
    title: string;
    series: string;
    mirrors?: string[];
  }

  type FrontMatterSpeakers = {
    [key: string]: string;
  }

  type FrontMatterCompletion = {
    timestamps: Boolean;
    notes: Boolean;
    issues: Boolean;
    mentions: Boolean;
    "content-verified": Boolean;
    content: Boolean;
    "speakers-identified": Boolean;
  }

  type FrontMatter = {
    source: FrontMatterSource;
    speakers: FrontMatterSpeakers;
    completion?: FrontMatterCompletion;
    added?: {
      author: string;
      date: string;
    }
  }

  type PagefindSubResultAnchor = {
    id: string;
  };

  type PagefindSubResult = {
    title: string;
    url: string;
    excerpt: string;
    anchor?: PagefindSubResultAnchor;
  };

  type PagefindFragment = {
    url: string;
    sub_results: PagefindSubResult[];
  };

  type PagefindResult = {
    id: string;
    score: number;
    words: string[];
    data: () => Promise<PagefindFragment>;
  };

  type PagefindResponse = {
    filters: PagefindFilters;
    results: PagefindResult[];
  }

  type ResultDataSection = {
    id: string;
    excerpt: string;
  };

  type ResultData = {
    slug: string;
    score: number;
    sections: ResultDataSection[];
  };

  interface HasExcerpt {
    excerpt: string | null;
  };

  type AsyncToken = Tokens.Generic & {
    resolveAsync?: (token: import("marked").Tokens.Generic, context: Context) => Promise<void>;
  };

  interface PagefindFilterValues extends Object {
    [value: string]: number;
  }

  interface PagefindFilters extends Object {
    [name: string]: PagefindFilterValues;
  };

  type Asset = {
    filename: string;
    date: string;
    slug: string;
    markdown: string;
    frontMatter: FrontMatter;
    html: string;
    issues: Issue[];
    sections: Section[];
    errors: string[];
    contributors: Contributor[];
  }

  type Context = {
    asset: Asset;
    fetcher: Fetcher;
    avatars: string[];
  };

  type Page = {
    asset: Asset;
    html: string;
    partial: string;
    filters: PagefindFilters;
  };

  type ThinAsset = {
    title: string;
    slug: string;
    date: string;
    publisher: string;
    sections: Section[];
    issues: Issue[];
  };
}

export { };

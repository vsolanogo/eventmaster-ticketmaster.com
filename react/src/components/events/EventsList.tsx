import React, { useEffect, useRef, useState } from "react";
import { useAppDispatch } from "../../store/store";
import {
  getEventsThunk,
  setEventsLimit,
  setInfiniteScroll,
  setPage,
  setSortBy,
  setSortOrder,
} from "../../redux/events/eventsSlice";
import {
  useAllEvents,
  useEventSortBy,
  useEventSortOrder,
  useEventsLimit,
  useEventsPage,
  useInfiniteScroll,
  useTotalCount,
} from "../../redux/selectors/selectorHooks";
import { Pagination, Spin, Switch } from "antd";
import { Link } from "wouter";
import { buttonTw } from "../../tailwind/tailwindClassNames";
import parse from "html-react-parser";

export const EventsList = () => {
  const dispatch = useAppDispatch();
  const totalCount = useTotalCount();
  const events = useAllEvents();
  const page = useEventsPage();
  const pageSize = useEventsLimit();
  const sortBy = useEventSortBy();
  const sortOrder = useEventSortOrder();
  const infiniteScroll = useInfiniteScroll();
  const divRef = useRef(null);
  const [isIntersecting, setIsIntersecting] = useState(false);

  useEffect(() => {
    dispatch(getEventsThunk(undefined));
  }, [sortBy, sortOrder, page, pageSize]);

  const handleSortChange = (e) => {
    dispatch(setSortBy(e.target.value));
  };

  const handleOrderChange = (e) => {
    dispatch(setSortOrder(e.target.value));
  };

  useEffect(() => {
    if (!infiniteScroll) {
      setPage(1);
    }
  }, [infiniteScroll]);

  useEffect(() => {
    const options = {
      root: null, // null means it will use the viewport
      rootMargin: "0px",
      threshold: 0.5,
    };

    const callback = (entries) => {
      entries.forEach((entry) => {
        setIsIntersecting(entry.isIntersecting);
      });
    };

    const observer = new IntersectionObserver(callback, options);
    if (divRef.current) {
      observer.observe(divRef.current);
    }

    return () => {
      if (divRef.current) {
        observer.unobserve(divRef.current);
      }
    };
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      if (infiniteScroll && totalCount > pageSize) {
        dispatch(setPage(1));
        dispatch(setEventsLimit(pageSize + 10));
      }
    }, 1000);

    return () => clearInterval(interval); // Clean up the interval
  }, [isIntersecting]);

  return (
    <>
      <div className="max-w-7xl mx-auto p-6">
        <div className="flex flex-col md:flex-row justify-end mb-4 items-center">
          <div className="mb-4 md:mb-0 md:mr-4">
            <label htmlFor="sort" className="mr-2">
              Sort by:
            </label>
            <select
              id="sort"
              className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:border-blue-300 pr-8"
              value={sortBy}
              onChange={handleSortChange}
            >
              <option value="eventDate">Event Date</option>
              <option value="title">Title</option>
              <option value="organizer">Organizer</option>
            </select>
          </div>

          <div className="mb-4 md:mb-0 md:mr-4">
            <label htmlFor="order" className="mr-2">
              Sort order:
            </label>
            <select
              id="order"
              className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:border-blue-300 pr-8"
              value={sortOrder}
              onChange={handleOrderChange}
            >
              <option value="ASC">Ascending</option>
              <option value="DESC">Descending</option>
            </select>
          </div>

          <div className="mb-4 md:mb-0 md:mr-4">
            <label className="mr-2">Infinite scroll</label>
            <Switch
              onChange={(e) => {
                dispatch(setInfiniteScroll(e));
              }}
            />
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {events.map((i) => (
            <div
              className="max-w-sm bg-white rounded-xl shadow-md overflow-hidden"
              key={i.id}
            >
              <div className="p-6">
                <h2 className="text-xl font-semibold mb-2">{i.title}</h2>
                <p className="text-gray-700 mb-4"> {parse(i.description)}</p>
                <div className="flex justify-between space-x-4">
                  <Link
                    href={`/events/register/${i.id}`}
                    className={`${buttonTw} mt-6 text-white`}
                  >
                    Register
                  </Link>

                  <Link
                    href={`/events/${i.id}`}
                    className={`${buttonTw} mt-6 text-white bg-green-500 text-white py-2 px-4 rounded`}
                  >
                    View
                  </Link>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {infiniteScroll && pageSize > totalCount && (
        <div className="text-center mt-4 text-gray-500">End of list</div>
      )}

      <div
        ref={divRef}
        className={"mt-10 mb-10"}
        style={{
          visibility:
            infiniteScroll && pageSize < totalCount ? "visible" : "hidden",
        }}
      >
        <Spin />
      </div>

      <Pagination
        onChange={(e) => {
          dispatch(setPage(e));
        }}
        onShowSizeChange={(current, size) => {
          dispatch(setEventsLimit(size));
        }}
        total={totalCount}
        defaultCurrent={1}
        defaultPageSize={10}
        current={page}
        pageSize={pageSize}
        disabled={infiniteScroll}
      />
    </>
  );
};

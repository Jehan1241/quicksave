import { useState, useEffect, useRef } from 'react'
import { useClickAway } from 'react-use'
import { useNavigate } from 'react-router-dom'
import { FaSortAmountUpAlt } from 'react-icons/fa'
import { IoIosArrowRoundBack, IoIosArrowRoundForward } from 'react-icons/io'
import { GoPlus } from 'react-icons/go'
import ImportPopUp from './ImportPopUp'

function NavBar(props) {
  const [tagsList, setTagsList] = useState()
  const [importClicked, setImportClicked] = useState(false)
  const [libraryClicked, setLibraryClicked] = useState(false)
  const [sortClicked, setSortClicked] = useState(false)
  const [filterClicked, setFilterClicked] = useState(false)
  const [order, setOrder] = useState(props.sortOrder)
  const navigate = useNavigate()
  const dropdownRef = useRef(null)

  const getAllTags = async () => {
    try {
      const response = await fetch(`http://localhost:8080/getAllTags`)
      const json = await response.json()
      setTagsList(json.tags)
    } catch (error) {
      console.error(error)
    }
  }
  useEffect(() => {
    getAllTags()
  }, [])

  useEffect(() => {
    setOrder(props.sortOrder)
  }, [props.sortOrder])

  /* Custom Hook from React-Use Library for Clickaway */
  useClickAway(dropdownRef, () => {
    setSortClicked(!sortClicked)
  })

  const sortClickHandler = () => {
    setSortClicked(!sortClicked)
  }

  const filterClickHandler = () => {
    setFilterClicked(!filterClicked)
  }

  const clearFilterClickHandler = async () => {
    try {
      const response = await fetch(`http://localhost:8080/clearFilter`)
      const json = await response.json()
      console.log(json)
    } catch (error) {
      console.error(error)
    }
  }

  const tagClickHandler = (tag) => {
    console.log(tag)
    setFilter(tag)
  }

  const setFilter = async (tag) => {
    try {
      const response = await fetch(`http://localhost:8080/setFilter?tag=${tag}`)
      const json = await response.json()
    } catch (error) {
      console.error(error)
    }
  }

  const sortOptionSelect = async (type) => {
    if (order == 'ASC') {
      setOrder('DESC')
      props.sortTypeChangeHandler(type, 'DESC')
    }
    if (order == 'DESC') {
      setOrder('ASC')
      props.sortTypeChangeHandler(type, 'ASC')
    }
  }

  const importClickHandler = () => {
    setImportClicked(!importClicked)
    setLibraryClicked(false)
    if (libraryClicked) {
      navigate('/')
    }
    console.log(importClicked)
  }

  const libraryClickHandler = () => {
    console.log('AA')
    setImportClicked(false)
    navigate('/')
  }

  /* In the next 2 funcs event.stopPropagation() needed to stop click through into next elements */
  const addGameManuallyClickHandler = (event) => {
    event.stopPropagation()
    setImportClicked(false)
    navigate('/AddGameManually')
  }

  const fromSteamClickHandler = (event) => {
    event.stopPropagation()
    setImportClicked(false)
    navigate('/AddGameSteam')
  }

  useEffect(() => {
    console.log('Current tileSize in NavBar:', props.tileSize) // Debugging log
  }, [props.tileSize])

  return (
    <>
      {/* Main BAR */}
      <div className="flex absolute z-30 flex-row justify-between items-center px-5 mx-10 mt-2 w-[calc(100vw-80px)] h-14 rounded-2xl shadow-md backdrop-blur-xl bg-primary/75">
        {/* Arrow Div */}
        <div className="flex justify-center items-center text-white">
          <button className="rounded-full hover:bg-gray-600/30">
            <IoIosArrowRoundBack size={40} onClick={() => navigate(-1)} />
          </button>
          <button className="rounded-full hover:bg-gray-600/30">
            <IoIosArrowRoundForward size={40} onClick={() => navigate(1)} />
          </button>
        </div>

        <div className="flex flex-row gap-2 justify-center items-center text-white">
          <input
            onChange={props.inputChangeHandler}
            className="px-3 my-auto h-7 rounded-xl bg-gray-600/20"
            placeholder="Search"
          />
          <div className="flex relative flex-col">
            <button>
              <FaSortAmountUpAlt
                size={18}
                onClick={sortClickHandler}
                className={`duration-150 ease-in-out ${sortClicked ? 'rotate-180' : ''}`}
              />
            </button>
            {sortClicked ? (
              <div
                ref={dropdownRef}
                className="flex absolute top-10 flex-col gap-2 p-1 text-sm rounded-lg bg-gameView"
              >
                <button
                  className="px-3 py-2 rounded-lg hover:bg-gray-600/30"
                  onClick={() => sortOptionSelect('CustomTitle')}
                >
                  Alphabetical {order == 'ASC' ? 'A Z' : 'Z A'}
                </button>

                <button
                  className="px-3 py-2 rounded-lg hover:bg-gray-600/30"
                  onClick={() => sortOptionSelect('CustomTimePlayed')}
                >
                  Time Played
                </button>

                <button
                  className="px-3 py-2 rounded-lg hover:bg-gray-600/30"
                  onClick={() => sortOptionSelect('CustomRating')}
                >
                  Rating
                </button>
              </div>
            ) : null}
          </div>
          <div>
            <button className="m-2" onClick={filterClickHandler}>
              filter
            </button>
            <button onClick={clearFilterClickHandler}>clear_filter</button>
            {filterClicked ? (
              <div className="flex overflow-y-scroll absolute top-14 flex-col gap-2 p-1 h-[75vh] text-sm rounded-lg bg-gameView">
                {tagsList.map((tag) => (
                  <button onClick={() => tagClickHandler(tag)}>{tag}</button>
                ))}
              </div>
            ) : null}
          </div>
        </div>
        {/* Add and Slider Div */}
        <div className="flex gap-3">
          <input
            value={props.tileSize}
            onChange={props.sizeChangeHandler}
            type="range"
            min={25}
            max={80}
          ></input>
          <button
            onClick={importClickHandler}
            className={`text-left text-white rounded-full hover:bg-gray-600/30`}
          >
            <GoPlus className={`inline p-2 rounded-full`} size={40} />
          </button>
        </div>
      </div>
      {/* POP BAR */}
      {importClicked ? (
        <ImportPopUp
          onGameAdded={props.onGameAdded}
          importClickHandler={importClickHandler}
          fromSteamClickHandler={fromSteamClickHandler}
          addGameManuallyClickHandler={addGameManuallyClickHandler}
        />
      ) : null}
    </>
  )
}

export default NavBar

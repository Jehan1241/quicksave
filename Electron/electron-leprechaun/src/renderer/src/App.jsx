import { useState, useEffect, useRef } from 'react'
import NavBar from './NavBar/NavBar'
import LibraryView from './LibraryView/LibraryView'
import { Route, Routes, useLocation } from 'react-router-dom'
import GameView from './GameView/GameView'

function App() {
  const [metaData, setMetaData] = useState([])
  const [searchText, setSearchText] = useState('')
  const location = useLocation()
  const state = location.state
  const [tileSize, setTileSize] = useState('default')
  const [debouncedTileSize, setDebouncedTileSize] = useState(tileSize)
  const [sortType, setSortType] = useState('default')
  const [sortOrder, setSortOrder] = useState('default')
  const [sse, setSse] = useState(null) // State to hold SSE connection
  const [timer, setTimer] = useState(null) // Timer state
  const [renderCounter, setRenderCounter] = useState(0) // New flag to check if it's the initial render
  const fetchDataRunningRef = useRef(false)

  const sortTypeChangeHandler = (type, order) => {
    setSortType(type)
    setSortOrder(order)
  }

  const NavBarInputChangeHandler = (e) => {
    const text = e.target.value
    setSearchText(text.toLowerCase())
  }

  const sizeChangeHandler = (e) => {
    const newSize = e.target.value
    setTileSize(newSize)
    // Clear the existing timer
    if (timer) {
      clearTimeout(timer)
    }
    // Set a new timer to update the debounced tile size
    const newTimer = setTimeout(() => {
      setDebouncedTileSize(newSize)
    }, 5000) // 5 seconds debounce
    // Update the timer state
    setTimer(newTimer)
  }

  useEffect(() => {
    fetchData()
    console.log('3')
    const eventSource = new EventSource('http://localhost:8080/sse-steam-updates')

    eventSource.onmessage = (event) => {
      console.log('SSE message received:', event.data)
      fetchData()
    }
    eventSource.onerror = (error) => {
      console.error('SSE Error:', error)
    }
    setSse(eventSource)
    return () => {
      eventSource.close()
    }
  }, [])

  const fetchData = async () => {
    console.log('Sending Get Basic Info')
    if (fetchDataRunningRef.current) return
    fetchDataRunningRef.current = true // Set flag to true
    try {
      const response = await fetch(
        `http://localhost:8080/getBasicInfo?type=${sortType}&order=${sortOrder}&size=${debouncedTileSize}`
      )
      const json = await response.json()
      console.log(json)
      setMetaData(json.MetaData)
      setSortOrder(json.SortOrder)
      setTileSize(json.Size)
    } catch (error) {
      console.error(error)
    } finally {
      fetchDataRunningRef.current = false // Reset the flag once fetch is done
    }
  }

  useEffect(() => {
    if (renderCounter < 2) {
      setRenderCounter(renderCounter + 1)
      console.log(renderCounter)
    } else {
      console.log('Here')
      fetchData()
    }
  }, [sortType, sortOrder, debouncedTileSize]) // Keep existing dependency for sortType and sortOrder

  const DataArray = Object.values(metaData)

  return (
    <div className="flex flex-col w-screen h-screen">
      <NavBar
        onGameAdded={fetchData}
        inputChangeHandler={NavBarInputChangeHandler}
        sizeChangeHandler={sizeChangeHandler}
        sortTypeChangeHandler={sortTypeChangeHandler}
        sortOrder={sortOrder}
        tileSize={tileSize}
      />
      <Routes>
        <Route
          element={<LibraryView tileSize={tileSize} searchText={searchText} data={DataArray} />}
          path="/"
        />
        <Route element={<GameView uid={state?.data} onDelete={fetchData} />} path="gameview" />
      </Routes>
    </div>
  )
}

export default App

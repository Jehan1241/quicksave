import { useState } from 'react'

function ListFoundGames(props) {
  const [commiting, setCommiting] = useState(false)

  const selectedGameClickHandler = async (appid) => {
    setCommiting(true)
    console.log(props.time)
    console.log(appid)
    console.log(props.SelectedPlatform)
    try {
      const response = await fetch('http://localhost:8080/InsertGameInDB', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          key: appid,
          platform: props.SelectedPlatform,
          time: props.time
        })
      })
      console.log(response)
      //props.onGameAdded()
      setCommiting(false)
    } catch (error) {
      console.error('Error:', error)
      setCommiting(false)
    }
  }

  if (props.FoundGames === '') {
    return
  } else {
    const data = JSON.parse(props.FoundGames.foundGames)
    if (Object.keys(data).length === 0) {
      return (
        <div className="text-left bg-gameView h-[28vh] w-auto overflow-scroll flex justify-center">
          No Games Found
        </div>
      )
    } else {
      return commiting ? (
        <div className="text-left bg-gameView h-[28vh] w-auto overflow-scroll flex justify-center">
          Adding Game...
        </div>
      ) : (
        <div className="overflow-scroll justify-center w-auto h-full text-left bg-gameView hover:border-gray-500">
          <div className="flex flex-col w-11/12">
            {Object.values(data).map((game) => (
              <button
                className="p-1 rounded-lg border-gray-700 bg-gameView hover:bg-gray-500/20"
                key={game.appid}
                onClick={() => selectedGameClickHandler(game.appid)}
              >
                <div className="flex justify-start text-left">{game.name}</div>
                <div className="flex justify-end">{new Date(game.date).getFullYear()}</div>
              </button>
            ))}
          </div>
        </div>
      )
    }
  }
}

export default ListFoundGames

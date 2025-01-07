import { useNavigate } from 'react-router-dom'

function Delete(props) {
  const navigate = useNavigate()

  const hardDeleteClickHandler = () => {
    console.log(`Delete ${props.uid}`)
    hardDelete()
  }

  const hideClickHandler = () => {
    console.log(`Hide ${props.uid}`)
    hide()
  }

  const hide = async () => {
    try {
      const response = await fetch(`http://localhost:8080/HideGame?uid=${props.uid}`)
      const json = await response.json()
    } catch (error) {
      console.error(error)
    }
    navigate('/', { replace: true })
  }

  const hardDelete = async () => {
    try {
      const response = await fetch(`http://localhost:8080/DeleteGame?uid=${props.uid}`)
      const json = await response.json()
    } catch (error) {
      console.error(error)
    }
    navigate('/', { replace: true })
  }

  return (
    <>
      <div className="flex flex-col gap-4 p-6 mt-4 w-full h-full text-base rounded-xl">
        <div className="flex flex-row gap-6 w-full h-full">
          <div className="flex flex-col gap-4 w-1/2">
            <div className="flex flex-row gap-4 items-center">
              <p>Hard Delete</p>
              <p className="text-sm">
                This will remove the game completely from the app, but future imports will bring it
                back
              </p>
              <button
                onClick={hardDeleteClickHandler}
                className="w-60 h-10 rounded-lg border-2 bg-gameView hover:bg-gray-500/20"
              >
                Hard Delete
              </button>
            </div>

            <div className="flex flex-row gap-4 items-center">
              <p>Hide</p>
              <p className="text-sm">
                This will remove the game from display and future imports will not bring it back
              </p>
              <button
                onClick={hideClickHandler}
                className="w-60 h-10 rounded-lg border-2 bg-gameView hover:bg-gray-500/20"
              >
                Hide
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}

export default Delete

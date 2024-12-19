function addScreenshotPopUp(props) {
  const addScreenshot = () => {
    console.log(document.getElementById('screenshotString').value)
    fetchData()
  }

  const fetchData = async () => {
    try {
      const response = await fetch(
        `http://localhost:8080/AddScreenshot?string=${document.getElementById('screenshotString').value}`
      )
      const json = await response.json()
      console.log(json)
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <div
      className="flex overflow-hidden fixed top-0 right-0 z-20 w-screen h-screen text-white bg-black/80"
      onClick={props.addScreenshotClickHandler}
    >
      <div
        className="flex flex-col p-5 m-auto w-2/3 h-2/3 text-2xl bg-gameView"
        onClick={(e) => e.stopPropagation()}
      >
        Screenshot Handler
        <div>
          <label>add ss</label>
          <input type="text" id="screenshotString" /> <br />
          <button onClick={addScreenshot}>add SS</button>
        </div>
      </div>
    </div>
  )
}

export default addScreenshotPopUp

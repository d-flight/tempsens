package data

type State struct {
	heatingState       HeatingState
	desiredTemperature Temperature
	latestReading      *Reading
	isUserControlled   bool
}

func NewState() *State {
	return &State{
		heatingState:       HEATING_STATE_IDLE,
		desiredTemperature: InvalidTemperature(),
		latestReading:      nil,
		isUserControlled:   false,
	}
}

func (s *State) GetHeatingState() HeatingState      { return s.heatingState }
func (s *State) SetHeatingState(state HeatingState) { s.heatingState = state }

func (s *State) GetLatestReading() *Reading        { return s.latestReading }
func (s *State) SetLatestReading(reading *Reading) { s.latestReading = reading }

func (s *State) GetDesiredTemperature() Temperature            { return s.desiredTemperature }
func (s *State) SetDesiredTemperature(tempearture Temperature) { s.desiredTemperature = tempearture }

func (s *State) IsUserControlled() bool                { return s.isUserControlled }
func (s *State) SetUserControlled(userControlled bool) { s.isUserControlled = userControlled }

package behavior

import (
	"time"

	"github.com/yenkeia/mirgo/game"
)

type Status uint8

const (
	SUCCESS Status = 1
	FAILED  Status = 2
	RUNNING Status = 3
	READY   Status = 4
)

type INode interface {
	Visit(*BT)
	Reset()
	Step()
	Status() Status
}

type BT struct {
	timer   time.Duration // 记录从启动开始的时间
	Root    INode
	Monster *game.Monster
}

func (c *BT) GetTime() time.Duration {
	return c.timer
}

func (c *BT) Process(dt time.Duration) {
	c.timer += dt
	c.Root.Visit(c)
	c.Root.Step()
}

func init() {
	game.SetMonsterBehaviorFactory(NewBehavior)
}

func NewBehavior(id int, mon *game.Monster) game.IBehavior {

	var root INode

	switch id {
	case 1, 2:
		root = DeerBrain()
	case 3:
		root = TreeBrain()
	case 4:
		root = SpittingSpiderBrain()
	case 6, 58:
		root = GuardBrain()
	case 57:
		root = TownArcherBrain() // ArcherGurad
	default:
		root = DefaultBrain()
	}

	/*
		case 1, 2:
			return new Deer(info);
		case 3:
			return new Tree(info);
		case 4:
			return new SpittingSpider(info);
		case 5:
			return new CannibalPlant(info);
		case 6:
			return new Guard(info);
		case 7:
			return new CaveMaggot(info);
		case 8:
			return new AxeSkeleton(info);
		case 9:
			return new HarvestMonster(info);
		case 10:
			return new FlamingWooma(info);
		case 11:
			return new WoomaTaurus(info);
		case 12:
			return new BugBagMaggot(info);
		case 13:
			return new RedMoonEvil(info);
		case 14:
			return new EvilCentipede(info);
		case 15:
			return new ZumaMonster(info);
		case 16:
			return new RedThunderZuma(info);
		case 17:
			return new ZumaTaurus(info);
		case 18:
			return new Shinsu(info);
		case 19:
			return new KingScorpion(info);
		case 20:
			return new DarkDevil(info);
		case 21:
			return new IncarnatedGhoul(info);
		case 22:
			return new IncarnatedZT(info);
		case 23:
			return new BoneFamiliar(info);
		case 24:
			return new DigOutZombie(info);
		case 25:
			return new RevivingZombie(info);
		case 26:
			return new ShamanZombie(info);
		case 27:
			return new Khazard(info);
		case 28:
			return new ToxicGhoul(info);
		case 29:
			return new BoneSpearman(info);
		case 30:
			return new BoneLord(info);
		case 31:
			return new RightGuard(info);
		case 32:
			return new LeftGuard(info);
		case 33:
			return new MinotaurKing(info);
		case 34:
			return new FrostTiger(info);
		case 35:
			return new SandWorm(info);
		case 36:
			return new Yimoogi(info);
		case 37:
			return new CrystalSpider(info);
		case 38:
			return new HolyDeva(info);
		case 39:
			return new RootSpider(info);
		case 40:
			return new BombSpider(info);
		case 41, 42:
			return new YinDevilNode(info);
		case 43:
			return new OmaKing(info);
		case 44:
			return new BlackFoxman(info);
		case 45:
			return new RedFoxman(info);
		case 46:
			return new WhiteFoxman(info);
		case 47:
			return new TrapRock(info);
		case 48:
			return new GuardianRock(info);
		case 49:
			return new ThunderElement(info);
		case 50:
			return new GreatFoxSpirit(info);
		case 51:
			return new HedgeKekTal(info);
		case 52:
			return new EvilMir(info);
		case 53:
			return new EvilMirBody(info);
		case 54:
			return new DragonStatue(info);
		case 55:
			return new HumanWizard(info);
		case 56:
			return new Trainer(info);
		case 57:
			return new TownArcher(info);
		case 58:
			return new Guard(info);
		case 59:
			return new HumanAssassin(info);
		case 60:
			return new VampireSpider(info);
		case 61:
			return new SpittingToad(info);
		case 62:
			return new SnakeTotem(info);
		case 63:
			return new CharmedSnake(info);
		case 64:
			return new IntelligentCreatureObject(info);
		case 65:
			return new MutatedManworm(info);
		case 66:
			return new CrazyManworm(info);
		case 67:
			return new DarkDevourer(info);
		case 68:
			return new Football(info);
		case 69:
			return new PoisonHugger(info);
		case 70:
			return new Hugger(info);
		case 71:
			return new Behemoth(info);
		case 72:
			return new FinialTurtle(info);
		case 73:
			return new TurtleKing(info);
		case 74:
			return new LightTurtle(info);
		case 75:
			return new WitchDoctor(info);
		case 76:
			return new HellSlasher(info);
		case 77:
			return new HellPirate(info);
		case 78:
			return new HellCannibal(info);
		case 79:
			return new HellKeeper(info);
		case 80:
			return new ConquestArcher(info);
		case 81:
			return new Gate(info);
		case 82:
			return new Wall(info);
		case 83:
			return new Tornado(info);
		case 84:
			return new WingedTigerLord(info);
		case 86:
			return new ManectricClaw(info);
		case 87:
			return new ManectricBlest(info);
		case 88:
			return new ManectricKing(info);
		case 89:
			return new IcePillar(info);
		case 90:
			return new TrollBomber(info);
		case 91:
			return new TrollKing(info);
		case 92:
			return new FlameSpear(info);
		case 93:
			return new FlameMage(info);
		case 94:
			return new FlameScythe(info);
		case 95:
			return new FlameAssassin(info);
		case 96:
			return new FlameQueen(info);
		case 97:
			return new HellKnight(info);
		case 98:
			return new HellLord(info);
		case 99:
			return new HellBomb(info);
		case 100:
			return new VenomSpider(info);
	*/

	bt := &BT{
		Root:    root,
		Monster: mon,
	}

	return bt
}

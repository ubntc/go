// Mode is an enum to control Thing object creation.
enum Mode {
    Default = 0,
    FirstNameOnly = 1,
    FirstNameAbbr = 2
}

// ThingOptions is an interface with optional properties.
// This is the only way to implement typed optional params with default values in TS.
interface ThingOptions {
    firstname?:string
    lastname?:string
    mode?:Mode
}

// Thing is a simple object with customizable properties.
class Thing {
    public id:string=""
    public name:string=""

    // constructor with optional customization options and default values.
    // I tried to add some complex constructor logic to show how you can use well-types options
    // to control object creation.
    constructor({firstname="Default", lastname="Thing", mode=Mode.Default}: ThingOptions) {
        this.id = (
            `${firstname.toLowerCase().replace(/[^a-z]+/g,'_')}` +
            `${lastname.toLocaleLowerCase().replace(/[^a-z]+/g,'_')}`
        )
        switch (mode) {
            case Mode.FirstNameOnly:  this.name = `${firstname}`;                 break;
            case Mode.FirstNameAbbr:  this.name = `${firstname[0]}. ${lastname}`; break;
            default:                  this.name = `${firstname} ${lastname}`
        }
    }
}

function main() {
    let t1 = new Thing({firstname:"The"})
    let t2 = new Thing({firstname:"Teddy", lastname:"Rex", mode:Mode.FirstNameAbbr})
    let t3 = new Thing({firstname:"What a thing!", mode:Mode.FirstNameOnly})
    console.log(t1,t2,t3)
}

main()

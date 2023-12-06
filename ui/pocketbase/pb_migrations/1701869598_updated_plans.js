/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("s7nhljkrmzzu8y6")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "0j9wrzc3",
    "name": "framework",
    "type": "relation",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "nwgspeyl1n5v30f",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("s7nhljkrmzzu8y6")

  // remove
  collection.schema.removeField("0j9wrzc3")

  return dao.saveCollection(collection)
})

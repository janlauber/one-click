/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("s7nhljkrmzzu8y6")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "9kslplki",
    "name": "frameworks",
    "type": "relation",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "nwgspeyl1n5v30f",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": null,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("s7nhljkrmzzu8y6")

  // remove
  collection.schema.removeField("9kslplki")

  return dao.saveCollection(collection)
})
